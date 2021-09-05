package gosql

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenSelect(es ...storage.TableNamer) {
	skipAuth := false

	// for _, p := range ps {
	// 	if p == "skip-auth" {
	// 		skipAuth = true
	// 	}
	// }

	ts := []string{}
	tms := []interface{}{}
	for _, e := range es {
		if !skipAuth {
			c.ReadAuthorization(e.TableName())
		}
		ts = append(ts, "%s.*")
		tms = append(tms, e.TableName())
	}

	q, vs := c.modifiers.Compile()
	c.values = append(c.values, vs...)

	c.query = fmt.Sprintf(strings.Join(ts, ","), tms...)
	c.query = fmt.Sprintf("SELECT %s FROM %s %s", c.query, es[0].TableName(), q)
}

func (c *Connection) ReadAuthorization(t string, params ...string) {
	perm := "r"

	permsZ := fmt.Sprintf("%%:%%:%%:%%%s%%", perm)
	permsA := fmt.Sprintf("%%:%%:%%%s%%:%%", perm)
	permsG := fmt.Sprintf("%%:%%%s%%:%%:%%", perm)
	permsU := fmt.Sprintf("%%%s%%:%%:%%:%%", perm)

	where := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsZ)

	if c.user != 0 {
		where.Or(fmt.Sprintf("%s.perms LIKE ?", t), permsA)

		if len(c.groups) > 0 {
			c.modifiers = append(c.modifiers, NewJoin(fmt.Sprintf("LEFT JOIN record_groups AS rgs_%s ON rgs_%[1]s.record_id = %[1]s.hash", t)))
			w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsG)
			groups := strings.Trim(strings.Replace(fmt.Sprint(c.groups), " ", ",", -1), "[]")
			w.And(fmt.Sprintf("rgs_%s.group_id IN (?)", t), groups)
			where.OrGroup(w)
		}

		c.modifiers = append(c.modifiers, NewJoin(fmt.Sprintf("LEFT JOIN record_users AS rus_%s ON rus_%[1]s.record_id = %[1]s.hash", t)))
		w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsU)
		w.And(fmt.Sprintf("rus_%s.user_id = ?", t), fmt.Sprint(c.user))
		where.OrGroup(w)

		where.Or(fmt.Sprintf("%s.owner_id = ?", t), c.user)
	}

	c.modifiers = append(c.modifiers, NewWhereGroup(where))
}

func (c *Connection) Read(es ...storage.Operator) {
	if c.err != nil {
		return
	}

	ts := []reflect.Type{}
	vs := []reflect.Value{}
	for _, e := range es {
		t := reflect.TypeOf(e)
		v := reflect.ValueOf(e)

		if t.Kind() == reflect.Ptr {
			t = t.Elem()
			if t.Kind() != reflect.Struct && t.Kind() != reflect.Slice {
				c.err = fmt.Errorf("not a struct")
				return
			}
			v = v.Elem()
		} else if t.Kind() == reflect.Map {
		} else {
			c.err = fmt.Errorf("not a pointer or map")
			return
		}

		ts = append(ts, t)
		vs = append(vs, v)
	}

	db, err := sql.Open(c.store.dbt, c.store.uri)
	if err != nil {
		c.err = fmt.Errorf("opening db connection, %w", err)
		return
	}
	defer func() {
		db.Close()
		c.modifiers = modifiers{}
		c.values = []interface{}{}
	}()
	db.SetConnMaxLifetime(3 * time.Second)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	tableNamers := make([]storage.TableNamer, len(es))
	for i := range es {
		tableNamers[i] = es[i]
	}
	c.GenSelect(tableNamers...)

	fmt.Printf("\n%s\n", c.query)
	rows, err := db.Query(c.query, c.values...)
	if err != nil {
		c.err = fmt.Errorf("reading from db, %w", err)
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		c.err = fmt.Errorf("getting row columns, %w", err)
		return
	}

	fmt.Println("cols", cols)

	nRows := 0
	for rows.Next() {
		tRow := ts[0]

		if ts[0].Kind() == reflect.Struct {
			values := []interface{}{}
			for i := range ts {
				values = append(values, dbToStruct(ts[i], vs[i])...)
			}
			err = rows.Scan(values...)
			if err != nil {
				c.err = chassis.Mark("scanning colunms", err)
			}
		} else if ts[0].Kind() == reflect.Map {
			values := []interface{}{}
			esRow := []storage.Operator{}
			vsRow := []reflect.Value{}
			for _, t := range ts {
				tRow = t.Elem()
				vRow := reflect.New(tRow).Elem()
				eRow, ok := vRow.Addr().Interface().(storage.Operator)
				if !ok {
					c.err = fmt.Errorf("not type storage.Operator")
					return
				}

				values = append(values, dbToStruct(tRow, vRow)...)
				esRow = append(esRow, eRow)
				vsRow = append(vsRow, vRow)
			}

			err = rows.Scan(values...)
			if err != nil {
				c.err = fmt.Errorf("scanning colunms, %w", err)
			}

			for i := range vs {
				if ts[i].Key().Kind() == reflect.Uint {
					key := reflect.ValueOf(esRow[i].IDValue())
					vs[i].SetMapIndex(key, vsRow[i])
				} else {
					key := reflect.ValueOf(esRow[i].UniqueCode())
					vs[i].SetMapIndex(key, vsRow[i])
				}
			}
		} else if ts[0].Kind() == reflect.Slice {
			values := []interface{}{}
			vsRow := []reflect.Value{}
			for _, t := range ts {
				tRow = t.Elem()
				vRow := reflect.New(tRow).Elem()

				values = append(values, dbToStruct(tRow, vRow)...)
				vsRow = append(vsRow, vRow)
			}

			err = rows.Scan(values...)
			if err != nil {
				c.err = fmt.Errorf("scanning colunms, %w", err)
			}

			for i := range vs {
				vs[i].Set(reflect.Append(vs[i], vsRow[i]))
			}
		}
		nRows++
	}
	fmt.Printf("read: %d, values: %v\n", nRows, c.values)
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

		if ft.Type.Kind() == reflect.Struct && ft.Type.Name() != "Time" {
			vs := dbToStruct(ft.Type, fv)
			values = append(values, vs...)
			continue
		}

		values = append(values, fv.Addr().Interface())
	}

	return values
}
