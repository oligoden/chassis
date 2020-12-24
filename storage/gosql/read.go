package gosql

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenSelect(e storage.TableNamer) {
	for _, m := range c.modifiers {
		q, vs := m.Compile()

		if strings.HasPrefix(q, " WHERE ") {
			q = strings.TrimPrefix(q, " WHERE ")
			w := NewWhere(q, vs...)
			if c.where == nil {
				c.where = NewWhereGroup(w)
			} else {
				c.where.AndGroup(w)
			}
		}

		if strings.HasPrefix(q, " LEFT JOIN ") {
			if c.join == nil {
				q = strings.TrimPrefix(q, " ")
				c.join = NewJoin(q)
			} else {
				qc, _ := c.join.Compile()
				qc = strings.TrimPrefix(qc, " ")
				c.join = NewJoin(qc + q)
			}
		}
	}

	c.ReadAuthorization(e.TableName())

	var q string
	if c.join != nil {
		j, _ := c.join.Compile()
		// j = strings.TrimPrefix(j, " ")
		q = q + j
	}

	if c.where != nil {
		where, vs := c.where.Compile()
		c.values = append(c.values, vs...)
		q = q + where
	}

	c.query = fmt.Sprintf("SELECT %s.* FROM %[1]s%s", e.TableName(), q)
}

func (c *Connection) ReadAuthorization(t string, params ...string) {
	perm := "r"

	for _, param := range params {
		if param == "with-update" {
			perm = "u"
		}
	}

	permsZ := fmt.Sprintf("%%:%%:%%:%%%s%%", perm)
	permsA := fmt.Sprintf("%%:%%:%%%s%%:%%", perm)
	permsG := fmt.Sprintf("%%:%%%s%%:%%:%%", perm)
	permsU := fmt.Sprintf("%%%s%%:%%:%%:%%", perm)

	where := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsZ)

	if c.user != 0 {
		where.Or(fmt.Sprintf("%s.perms LIKE ?", t), permsA)

		if c.join == nil {
			c.join = NewJoin(fmt.Sprintf("LEFT JOIN record_groups on record_groups.record_id = %s.hash", t))
		} else {
			jc, _ := c.join.Compile()
			c.join = NewJoin(jc + fmt.Sprintf(" LEFT JOIN record_groups on record_groups.record_id = %s.hash", t))
		}

		if len(c.groups) > 0 {
			w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsG)
			groups := strings.Trim(strings.Replace(fmt.Sprint(c.groups), " ", ",", -1), "[]")
			w.And("record_groups.group_id IN (?)", groups)
			where.OrGroup(w)
		}

		c.join.Add(fmt.Sprintf("LEFT JOIN record_users on record_users.record_id = %s.hash", t))
		w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsU)
		w.And("record_users.user_id = ?", fmt.Sprint(c.user))
		where.OrGroup(w)

		where.Or(fmt.Sprintf("%s.owner_id = ?", t), c.user)
	}

	if c.where == nil {
		c.where = NewWhereGroup(where)
	} else {
		c.where.AndGroup(where)
	}
}

func (c *Connection) Read(e storage.Operator) {
	if c.err != nil {
		return
	}

	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() != reflect.Struct && t.Kind() != reflect.Slice {
			c.err = fmt.Errorf("not a struct")
			return
		}
		v = v.Elem()
	} else if t.Kind() == reflect.Slice {
	} else if t.Kind() == reflect.Map {
	} else {
		c.err = fmt.Errorf("not a pointer, map or slice")
		return
	}

	db, err := sql.Open(c.store.dbt, c.store.uri)
	if err != nil {
		c.err = fmt.Errorf("opening db connection, %w", err)
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.GenSelect(e)
	log.Println(c.query, c.values)

	rows, err := db.Query(c.query, c.values...)
	if err != nil {
		c.err = fmt.Errorf("reading from db, %w", err)
		return
	}

	c.where = nil
	c.join = nil

	// cols, err := rows.Columns()
	// if err != nil {
	// 	c.err = fmt.Errorf("getting row columns, %w", err)
	// 	return
	// }

	// fmt.Println("cols", cols)

	for rows.Next() {
		tRow := t

		if t.Kind() == reflect.Struct {
			values := dbToStruct(t, v)
			err = rows.Scan(values...)
			if err != nil {
				c.err = fmt.Errorf("scanning colunms, %w", err)
			}
		} else if t.Kind() == reflect.Map {
			tRow = t.Elem()
			vRow := reflect.New(tRow).Elem()
			eRow, ok := vRow.Addr().Interface().(storage.Operator)
			if !ok {
				c.err = fmt.Errorf("not type storage.Storer")
				return
			}

			values := dbToStruct(tRow, vRow)
			err = rows.Scan(values...)
			if err != nil {
				c.err = fmt.Errorf("scanning colunms, %w", err)
			}

			vUC := reflect.ValueOf(eRow.UniqueCode())
			v.SetMapIndex(vUC, vRow)
		} else if t.Kind() == reflect.Slice {
			tRow = t.Elem()
			vRow := reflect.New(tRow).Elem()

			values := dbToStruct(tRow, vRow)
			err = rows.Scan(values...)
			if err != nil {
				c.err = fmt.Errorf("scanning colunms, %w", err)
			}
			v.Set(reflect.Append(v, vRow))
		}
	}
}

func dbToStruct(t reflect.Type, v reflect.Value) []interface{} {
	values := []interface{}{}

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		// fmt.Printf("%d. %v (%v, %v), tag: '%v' canset %v\n", i+1, ft.Name, ft.Type.Name(), ft.Type.Kind(), ft.Tag.Get("gosql"), fv.CanSet())

		if tag, got := ft.Tag.Lookup("gosql"); got {
			if tag == "-" {
				continue
			}
		}

		if ft.PkgPath != "" {
			continue
		}

		if ft.Type.Kind() == reflect.Struct {
			vs := dbToStruct(ft.Type, fv)
			values = append(values, vs...)
			continue
		}

		values = append(values, fv.Addr().Interface())
	}

	return values
}