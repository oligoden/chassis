package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenSelect(e storage.TableNamer) {
	var q string

	c.ReadAuthorization(e.TableName())

	if c.join != nil {
		q = q + c.join.Compile()
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

		c.join = NewJoin(fmt.Sprintf("LEFT JOIN record_groups on record_groups.record_id = %s.hash", t))
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

	c.Where(where)
}

func (c *Connection) Read(e storage.Storer) {
	if c.store.err != nil {
		return
	}

	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			c.store.err = fmt.Errorf("not a struct")
			return
		}
		v = v.Elem()
	} else if t.Kind() == reflect.Slice {
	} else if t.Kind() == reflect.Map {
	} else {
		c.store.err = fmt.Errorf("not a pointer, map or slice")
		return
	}

	db, err := sql.Open(c.store.dbt, c.store.uri)
	if err != nil {
		c.store.err = fmt.Errorf("opening db connection, %w", err)
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(10 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	c.GenSelect(e)

	rows, err := db.Query(c.query, c.values...)
	if err != nil {
		c.store.err = fmt.Errorf("reading from db, %w", err)
		return
	}

	c.where = nil
	c.join = nil

	// cols, err := rows.Columns()
	// if err != nil {
	// 	c.store.err = fmt.Errorf("getting row columns, %w", err)
	// 	return
	// }

	// fmt.Println("cols", cols)

	for rows.Next() {
		tRow := t

		if t.Kind() == reflect.Map {
			tRow = t.Elem()
			v = reflect.New(tRow).Elem()
		}

		values, uc := dbToStruct(tRow, v)
		err = rows.Scan(values...)
		if err != nil {
			c.store.err = fmt.Errorf("scanning colunms, %w", err)
		}

		if t.Kind() == reflect.Map {
			reflect.ValueOf(e).SetMapIndex(uc, v)
		}
	}
}

func dbToStruct(t reflect.Type, v reflect.Value) ([]interface{}, reflect.Value) {
	values := []interface{}{}
	var uc reflect.Value

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		if ft.Name == "UC" {
			uc = fv
		}

		// fmt.Printf("%d. %v (%v, %v), tag: '%v' canset %v\n", i+1, ft.Name, ft.Type.Name(), ft.Type.Kind(), ft.Tag.Get("gosql"), fv.CanSet())

		if tag, got := ft.Tag.Lookup("gosql"); got {
			if tag == "-" {
				continue
			}
		}

		if ft.PkgPath != "" {
			continue
		}

		values = append(values, fv.Addr().Interface())
	}

	return values, uc
}
