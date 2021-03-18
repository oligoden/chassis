package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenUpdate(e storage.Operator) {
	q := ""

	vs, err := structToUpdateQ(e, &q)
	if err != nil {
		c.Err(err)
	}

	c.values = []interface{}{}
	c.values = append(c.values, vs...)

	wq, wvs := c.modifiers.Compile()
	if !strings.Contains(wq, "WHERE") {
		w := NewWhere("id = ?", e.IDValue())
		c.AddModifiers(w)
		wq, wvs = c.modifiers.Compile()
	}
	c.values = append(c.values, wvs...)

	c.query = fmt.Sprintf("UPDATE %s SET %s %s", e.TableName(), q, wq)
}

func structToUpdateQ(e interface{}, q *string) ([]interface{}, error) {
	values := []interface{}{}
	sep := ""
	if len(*q) > 0 {
		sep = ", "
	}

	t := reflect.TypeOf(e)
	v := reflect.ValueOf(e)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() != reflect.Struct {
			return []interface{}{}, fmt.Errorf("not a struct")
		}
		v = v.Elem()
	} else if t.Kind() == reflect.Slice {
	} else if t.Kind() == reflect.Map {
	} else {
		return []interface{}{}, fmt.Errorf("not a pointer, map or slice")
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == "ID" {
			continue
		}

		if t.Field(i).Name == "UC" {
			continue
		}

		if t.Field(i).Name == "OwnerID" {
			continue
		}

		if t.Field(i).Name == "Perms" {
			continue
		}

		ft := t.Field(i)
		fv := v.Field(i)

		if tag, got := ft.Tag.Lookup("gosql"); got {
			if tag == "-" {
				continue
			}
		}

		if ft.PkgPath != "" {
			continue
		}

		if ft.Type.Kind() == reflect.Struct && ft.Type.Name() != "Time" {
			vs, err := structToUpdateQ(fv.Addr().Interface(), q)
			if err != nil {
				return []interface{}{}, err
			}
			values = append(values, vs...)
			continue
		}

		values = append(values, fv.Interface())

		*q = *q + sep + ToSnakeCase(ft.Name) + " = ?"
		sep = ", "
	}

	return values, nil
}

func (c *Connection) Update(e storage.Operator) {
	if c.err != nil {
		return
	}

	if auth, err := Authorize(e, "u", c.user, c.groups); !auth {
		if err != nil {
			c.err = err
			return
		}
		c.err = errors.New("update authorization failed")
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

	c.GenUpdate(e)

	result, err := db.Exec(c.query, c.values...)
	if err != nil {
		c.err = err
	}

	updated := int64(0)
	if result != nil {
		updated, _ = result.RowsAffected()
	}

	fmt.Printf("\n%s\nupdated: %d, values: %v\n", c.query, updated, c.values)
}
