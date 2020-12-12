package gosql

import (
	"database/sql"

	"github.com/oligoden/chassis/storage"
)

type Connection struct {
	store     *Store
	modifiers []storage.Modifier
	join      *Join
	where     *Where
	query     string
	values    []interface{}
	user      uint
	groups    []uint
	db        *sql.DB
	logger    Logger
	err       error
}

func NewConnection(user uint, groups []uint) *Connection {
	return &Connection{
		user:   user,
		groups: groups,
	}
}

func (s *Store) Connect(user uint, groups []uint) storage.Crudder {
	return &Connection{
		store:  s,
		user:   user,
		groups: groups,
	}
}

func (c Connection) Query() (string, []interface{}) {
	return c.query, c.values
}

func (c *Connection) AddModifiers(m ...storage.Modifier) {
	c.modifiers = append(c.modifiers, m...)
}

func (c *Connection) Err(e ...error) error {
	if len(e) > 0 {
		c.err = e[0]
	}
	return c.err
}

type Logger interface {
	Log(interface{})
}

func (c *Connection) SetLogger(l Logger) {
	c.logger = l
}
