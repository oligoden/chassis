package gosql

import (
	"database/sql"
	"sort"
	"strings"

	"github.com/oligoden/chassis/storage"
)

type Connection struct {
	store     *Store
	modifiers modifiers
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

type modifiers []storage.Modifier

func (ms modifiers) Len() int {
	return len(ms)
}

func (ms modifiers) Swap(i, j int) {
	ms[i], ms[j] = ms[j], ms[i]
}

func (ms modifiers) Less(i, j int) bool {
	return ms[i].Order() < ms[j].Order()
}

func (ms modifiers) Compile() (string, []interface{}) {
	sort.Sort(ms)

	var qs []string
	var vsAll []interface{}

	for i, m := range ms {
		var q string
		var vs []interface{}

		if i > 0 {
			if ms[i-1].Order() != m.Order() {
				q, vs = m.Compile("first")
			} else {
				q, vs = m.Compile("same")
			}
		} else {
			q, vs = m.Compile("first")
		}

		qs = append(qs, q)
		vsAll = append(vsAll, vs...)
	}
	return strings.Join(qs, " "), vsAll
}
