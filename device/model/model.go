package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
	"github.com/oligoden/chassis/storage/gosql"
)

type Operator interface {
	Create()
	Read()
	Update()
	Delete()
	Communicator
	DataSelector
}

type Communicator interface {
	Bind()
	User() (uint, []uint)
	Session() uint
	Hasher()
	Err(...interface{}) error
}

type DataSelector interface {
	Data(...data.Operator) data.Operator
	NewData(...string)
}

type Default struct {
	Request *http.Request `json:"-"`
	Hash    string        `json:"hash"`
	sesh    uint
	user    uint
	groups  []uint
	err     []error
	urlbase string
	data    data.Operator
	Store   Connector
}

type Connector interface {
	Connect(user uint, groups []uint) storage.Crudder
	Rnd() *rand.Rand
}

func (m Default) User() (uint, []uint) {
	return m.user, m.groups
}

func (m Default) Session() uint {
	return m.sesh
}

func (m *Default) SetBaseURL(base string) {
	m.urlbase = base
}

func (m *Default) BindUser(usg ...uint) {
	if m.Err() != nil {
		return
	}

	if len(usg) == 0 && m.Request == nil {
		fmt.Println("no user set")
		return
	}

	if len(usg) >= 1 {
		m.user = usg[0]
	} else {
		u := m.Request.Header.Get("X_user")
		if u == "" {
			m.Err(chassis.Mark("X_user not set, expected >= 0"))
		}
		user, err := strconv.Atoi(u)
		if err != nil {
			m.Err(chassis.Mark("user binding X_user", err))
			return
		}
		m.user = uint(user)
	}

	if len(usg) >= 2 {
		m.sesh = usg[1]
	} else {
		s := m.Request.Header.Get("X_session")
		if s == "" {
			m.Err(chassis.Mark("X_session not set"))
		}
		sesh, err := strconv.Atoi(s)
		if err != nil {
			m.Err(chassis.Mark("session binding X_session", err))
			return
		}
		m.sesh = uint(sesh)
	}

	if len(usg) >= 3 {
		m.groups = append(m.groups, usg[2:]...)
	}

	if m.Request == nil {
		return
	}

	if m.Request.Header.Get("X_user_groups") != "" {
		for _, g := range strings.Split(m.Request.Header.Get("X_user_groups"), ",") {
			group, err := strconv.Atoi(g)
			if err != nil {
				m.Err(chassis.Mark("user binding X_user_groups", err))
				return
			}
			m.groups = append(m.groups, uint(group))
		}
	}
}

func (m *Default) Bind() {
	if m.Err() != nil {
		return
	}

	if m.Request == nil {
		m.Err("request not set")
		return
	}

	if m.data == nil {
		m.Err("no data set")
		return
	}

	err := m.bind()
	if err != nil {
		m.Err(err)
		return
	}
}

func (m *Default) Data(d ...data.Operator) data.Operator {
	if len(d) > 0 {
		m.data = d[0]
	}
	return m.data
}

func (m *Default) NewData(ds ...string) {}

func (m *Default) Hasher() {
	json, err := json.Marshal(m.data)
	if err != nil {
		m.Err(err)
		return
	}
	h := sha1.New()
	h.Write(json)
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Default) Err(es ...interface{}) error {
	if m.err == nil {
		m.err = []error{}
	}

	for _, e := range es {
		switch t := e.(type) {
		case string:
			m.err = append(m.err, errors.New(t))
		case error:
			m.err = append(m.err, t)
		default:
			m.err = append(m.err, errors.New("unknown error type"))
		}
	}

	if len(m.err) == 0 {
		return nil
	}

	return m.err[0]
}

func (m *Default) Read() {
	if m.Err() != nil {
		return
	}

	c := m.Store.Connect(m.User())

	if m.data.UniqueCode() != "" {
		w := gosql.NewWhere("uc = ?", m.data.UniqueCode())
		c.AddModifiers(w)
	}
	c.Read(m.data)

	err := c.Err()
	if err != nil {
		m.Err(err)
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Err(err)
		return
	}

	m.Hasher()
}

func (m *Default) Delete() {
	if m.Err() != nil {
		return
	}

	c := m.Store.Connect(m.User())

	if m.data.UniqueCode() != "" {
		where := gosql.NewWhere("uc = ?", m.data.UniqueCode())
		c.AddModifiers(where)
	}

	c.Delete(m.data)
	err := c.Err()
	if err != nil {
		m.Err(err)
		return
	}
}
