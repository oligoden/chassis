package gosql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/oligoden/chassis/storage"
)

func (c *Connection) GenDelete(e storage.TableNamer) {
	c.DeleteAuthorization(e.TableName())

	q, vs := c.modifiers.Compile()
	c.values = append(c.values, vs...)

	c.query = fmt.Sprintf("DELETE %s.* FROM %[1]s %s", e.TableName(), q)
}

func (c *Connection) DeleteAuthorization(t string, params ...string) {
	perm := "d"

	permsZ := fmt.Sprintf("%%:%%:%%:%%%s%%", perm)
	permsA := fmt.Sprintf("%%:%%:%%%s%%:%%", perm)
	permsG := fmt.Sprintf("%%:%%%s%%:%%:%%", perm)
	permsU := fmt.Sprintf("%%%s%%:%%:%%:%%", perm)

	where := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsZ)

	if c.user != 0 {
		where.Or(fmt.Sprintf("%s.perms LIKE ?", t), permsA)

		c.modifiers = append(c.modifiers, NewJoin(fmt.Sprintf("LEFT JOIN record_groups on record_groups.record_id = %s.hash", t)))

		if len(c.groups) > 0 {
			w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsG)
			groups := strings.Trim(strings.Replace(fmt.Sprint(c.groups), " ", ",", -1), "[]")
			w.And("record_groups.group_id IN (?)", groups)
			where.OrGroup(w)
		}

		c.modifiers = append(c.modifiers, NewJoin(fmt.Sprintf("LEFT JOIN record_users on record_users.record_id = %s.hash", t)))
		w := NewWhere(fmt.Sprintf("%s.perms LIKE ?", t), permsU)
		w.And("record_users.user_id = ?", fmt.Sprint(c.user))
		where.OrGroup(w)

		where.Or(fmt.Sprintf("%s.owner_id = ?", t), c.user)
	}

	c.modifiers = append(c.modifiers, NewWhereGroup(where))
}

func (c *Connection) Delete(e storage.Operator) {
	if c.err != nil {
		return
	}

	db, err := sql.Open(c.store.dbt, c.store.uri)
	if err != nil {
		c.err = fmt.Errorf("opening db connection, %w", err)
		return
	}
	defer db.Close()
	db.SetConnMaxLifetime(3 * time.Second)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	c.GenDelete(e)

	result, err := db.Exec(c.query, c.values...)
	if err != nil {
		c.err = fmt.Errorf("reading from db, %w", err)
	}

	deleted := int64(0)
	if result != nil {
		deleted, _ = result.RowsAffected()
	}
	fmt.Printf("\n%s\ndeleted: %d, values: %v\n", c.query, deleted, c.values)

	c.modifiers = modifiers{}
	c.values = []interface{}{}
}
