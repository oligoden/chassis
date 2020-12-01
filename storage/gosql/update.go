package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenUpdate(incoming, existing storage.Storer) {
	tExisting := reflect.TypeOf(existing)
	vExisting := reflect.ValueOf(existing)

	if tExisting.Kind() == reflect.Ptr {
		tExisting = tExisting.Elem()
		if tExisting.Kind() != reflect.Struct {
			c.store.err = fmt.Errorf("not a struct")
			return
		}
		vExisting = vExisting.Elem()
	} else {
		c.store.err = fmt.Errorf("not a pointer")
		return
	}

	tIncomming := reflect.TypeOf(incoming)
	vIncomming := reflect.ValueOf(incoming)

	if tIncomming.Kind() == reflect.Ptr {
		tIncomming = tIncomming.Elem()
		if tIncomming.Kind() != reflect.Struct {
			c.store.err = fmt.Errorf("not a struct")
			return
		}
		vIncomming = vIncomming.Elem()
	} else {
		c.store.err = fmt.Errorf("not a pointer")
		return
	}

	q := ""
	var idValue reflect.Value
	c.values = []interface{}{}
	sep := ""
	for i := 0; i < tExisting.NumField(); i++ {
		if tExisting.Field(i).Name == "ID" {
			idValue = vExisting.Field(i)
			continue
		}

		if tExisting.Field(i).Name == "UC" {
			continue
		}

		if tExisting.Field(i).Name == "Hash" {
			continue
		}

		if tag, got := tIncomming.Field(i).Tag.Lookup("gosql"); got {
			if tag == "-" {
				continue
			}
		}

		if tIncomming.Field(i).PkgPath != "" {
			continue
		}

		if !reflect.DeepEqual(vExisting.Field(i).Interface(), vIncomming.Field(i).Interface()) {
			q = q + sep + ToSnakeCase(tIncomming.Field(i).Name) + " = ?"
			c.values = append(c.values, vIncomming.Field(i).Interface())
			sep = ", "
		}
	}

	c.values = append(c.values, idValue.Interface())
	c.query = fmt.Sprintf("UPDATE %s SET %s WHERE id = ?", existing.TableName(), q)
}

func (c *Connection) Update(eIncoming, eExisting storage.Storer) {
	if c.store.err != nil {
		return
	}

	if auth, err := Authorize(eExisting, "u", c.user, c.groups); !auth {
		if err != nil {
			c.store.err = err
			return
		}
		c.store.err = errors.New("create authorization failed")
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

	c.GenUpdate(eIncoming, eExisting)
	fmt.Println(c.query, c.values)

	_, err = db.Exec(c.query, c.values...)
	if err != nil {
		c.store.err = err
	}
}
