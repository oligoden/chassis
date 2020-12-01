package gosql

import (
	"database/sql"
)

type Connection struct {
	store  *Store
	join   *Join
	where  *Where
	query  string
	values []interface{}
	user   uint
	groups []uint
	db     *sql.DB
}

func (s *Store) Connect(user uint, groups []uint) *Connection {
	return &Connection{
		store:  s,
		user:   user,
		groups: groups,
	}
}

func (c Connection) Query() (string, []interface{}) {
	return c.query, c.values
}

func (c *Connection) Where(w *Where) {
	if c.where == nil {
		c.where = NewWhereGroup(w)
	} else {
		c.where.AndGroup(w)
	}
}
