package gosql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenInsert(e storage.TableNamer) {
	// tablename := inflection.Plural(m.TableName())
	// tablename = strings.ToLower(tablename)

	q1 := ""
	q2 := ""

	t := reflect.TypeOf(e).Elem()
	v := reflect.ValueOf(e).Elem()

	if t.Kind() != reflect.Struct {
		c.store.err = fmt.Errorf("not a struct")
		return
	}

	c.values = []interface{}{}
	sep := ""
	for i := 0; i < t.NumField(); i++ {
		if i > 0 {
			sep = ", "
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

		c.values = append(c.values, fv.Interface())

		q1 = q1 + sep + ToSnakeCase(ft.Name)
		q2 = q2 + sep + "?"

		// fmt.Printf("%d. %v (%v, %v), tag: '%v'\n", i+1, ft.Name, ft.Type.Name(), ft.Type.Kind(), ft.Tag.Get("form"))

		// if ft.Type.Name() == "RecordDefault" {
		// 	err := m.structBind(fv.Addr().Interface())
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// if tag, got := ft.Tag.Lookup("form"); got {
		// 	if val := m.Request.FormValue(tag); val != "" {
		// 		setType(ft.Type.Kind(), val, fv)
		// 	}
		// }
	}

	c.query = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", e.TableName(), q1, q2)
}

func (c *Connection) Create(e storage.Storer) {
	if c.store.err != nil {
		return
	}

	e.Owner(c.user)

	if auth, err := Authorize(e, "c", c.user, c.groups); !auth {
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
	c.db = db

	e.UniqueCode(c.store.UniqueCodeFunc()(c.store.UniqueCodeLength()))
	c.GenInsert(e)

	_, err = c.db.Exec(c.query, c.values...)
	if err != nil {
		if isDuplicateUC(err) {
			c.store.err = c.retryCreate(e)
		} else {
			c.store.err = err
		}
	}
}

func isDuplicateUC(err error) bool {
	return strings.Contains(err.Error(), "Error 1062: Duplicate entry") && strings.Contains(err.Error(), "for key 'uc'")
}

func (c *Connection) retryCreate(e storage.Storer) error {
	for i := 0; i < 3; i++ {
		e.UniqueCode(c.store.UniqueCodeFunc()(c.store.UniqueCodeLength()))
		c.GenInsert(e)

		_, err := c.db.Exec(c.query, c.values...)
		if err != nil {
			if !isDuplicateUC(err) {
				return c.retryCreate(e)
			}
		} else {
			return nil
		}
	}

	c.store.uniqueCodeLength++
	e.UniqueCode(c.store.UniqueCodeFunc()(c.store.uniqueCodeLength))
	c.GenInsert(e)

	_, err := c.db.Exec(c.query, c.values...)
	if err != nil {
		return err
	}
	return nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
