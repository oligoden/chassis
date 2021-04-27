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
	q1 := ""
	q2 := ""

	err, vs := reflectStruct(e, &q1, &q2)
	if err != nil {
		c.Err(err)
	}

	c.query = fmt.Sprintf("INSERT INTO %s(%s) VALUES(%s)", e.TableName(), q1, q2)
	c.values = vs
}

func reflectStruct(e interface{}, q1, q2 *string) (error, []interface{}) {
	values := []interface{}{}
	sep := ""
	if len(*q1) > 0 {
		sep = ", "
	}

	t := reflect.TypeOf(e).Elem()
	v := reflect.ValueOf(e).Elem()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct"), []interface{}{}
	}

	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).Name == "ID" {
			continue
		}

		if t.Field(i).Name == "TS" {
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
			err, vs := reflectStruct(fv.Addr().Interface(), q1, q2)
			if err != nil {
				return err, []interface{}{}
			}
			values = append(values, vs...)
			continue
		}

		values = append(values, fv.Interface())

		*q1 = *q1 + sep + ToSnakeCase(ft.Name)
		*q2 = *q2 + sep + "?"
		sep = ", "
	}

	return nil, values
}

func (c *Connection) Create(e storage.Operator) {
	if c.err != nil {
		return
	}

	e.Owner(c.user)

	if auth, reason, err := Authorize(e, "c", c.user, c.groups); !auth {
		if err != nil {
			c.err = err
			return
		}
		c.err = errors.New("create authorization failed, " + reason)
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
	c.db = db

	e.UniqueCode(c.store.UniqueCodeFunc()(c.store.UniqueCodeLength()))
	c.GenInsert(e)

	if c.logger != nil {
		c.logger.Log("")
	}

	fmt.Printf("\n%s\nvalues: %v\n", c.query, c.values)
	result, err := c.db.Exec(c.query, c.values...)
	if err != nil {
		if isDuplicateUC(err) {
			result, err = c.retryCreate(e)
			if err != nil {
				c.err = err
				return
			}
		} else {
			c.err = err
			return
		}
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.err = err
		return
	}

	created := int64(0)
	if result != nil {
		created, _ = result.RowsAffected()
	}

	fmt.Printf("created: %d\n", created)
	e.IDValue(uint(id))
}

func isDuplicateUC(err error) bool {
	return strings.Contains(err.Error(), "Error 1062: Duplicate entry") && strings.Contains(err.Error(), "for key 'uc'")
}

func (c *Connection) retryCreate(e storage.Operator) (sql.Result, error) {
	for i := 0; i < 3; i++ {
		e.UniqueCode(c.store.UniqueCodeFunc()(c.store.UniqueCodeLength()))
		c.GenInsert(e)

		r, err := c.db.Exec(c.query, c.values...)
		if err != nil {
			if isDuplicateUC(err) {
				continue
			}
		}
		return r, err
	}

	c.store.uniqueCodeLength++
	e.UniqueCode(c.store.UniqueCodeFunc()(c.store.uniqueCodeLength))
	c.GenInsert(e)

	r, err := c.db.Exec(c.query, c.values...)
	return r, err
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
