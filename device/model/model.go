package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
)

type Operator interface {
	// Manage(storage.DBManager, string)
	Create()
	Read()
	Update()
	// Append(string, storage.DBReader)
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
}

type Default struct {
	Request *http.Request        `json:"-"`
	NewData func() data.Operator `json:"-"`
	Hash    string               `json:"hash"`
	sesh    uint
	user    uint
	groups  []uint
	err     []error
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

func (m *Default) BindUser() {
	if m.Err() != nil {
		return
	}

	if m.Request == nil {
		log.Println("request not set")
		return
	}

	u := m.Request.Header.Get("X_user")
	user, err := strconv.Atoi(u)
	if err != nil {
		m.Err(fmt.Errorf("user binding X_user, %w", err))
		return
	}
	m.user = uint(user)

	s := m.Request.Header.Get("X_session")
	sesh, err := strconv.Atoi(s)
	if err != nil {
		m.Err(fmt.Errorf("session binding X_session, %w", err))
		return
	}
	m.sesh = uint(sesh)

	if m.Request.Header.Get("X_user_groups") != "" {
		for _, g := range strings.Split(m.Request.Header.Get("X_user_groups"), ",") {
			group, err := strconv.Atoi(g)
			if err != nil {
				m.Err(fmt.Errorf("user binding X_user_groups, %w", err))
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

// func (m *Default) Manage(db storage.DBManager, action string) {
// 	if m.Err() != nil {
// 		return
// 	}

// 	db.Manage(m.data, action)
// 	err := db.Error()
// 	if err != nil {
// 		m.Err(err)
// 		return
// 	}
// }

func (m *Default) Create() {
	if m.Err() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Err(err)
		return
	}

	c := m.Store.Connect(m.User())
	c.Create(m.data)
	err = c.Err()
	if err != nil {
		m.Err(err)
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.Err(err)
		return
	}

	c.Update(m.data)
	err = c.Err()
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

func (m *Default) Read() {
	if m.Err() != nil {
		return
	}

	c := m.Store.Connect(m.User())
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

func (m *Default) Update() {
	if m.Err() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Err(err)
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.Err(err)
		return
	}

	c := m.Store.Connect(m.User())
	c.Update(m.data)
	err = c.Err()
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
