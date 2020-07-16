package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
)

type Operator interface {
	Manage(storage.DBManager, string)
	Create(storage.DBCreater)
	Read(storage.DBReader, ...string)
	Update(storage.DBUpdater)
	// Append(string, storage.DBReader)
	Communicator
	DataSelector
}

type Communicator interface {
	Bind()
	User() (uint, []uint)
	Hasher()
	Error(...interface{}) error
}

type DataSelector interface {
	Data(...data.Operator) data.Operator
}

type Default struct {
	Request *http.Request        `json:"-"`
	NewData func() data.Operator `json:"-"`
	Hash    string               `json:"hash"`
	user    uint
	groups  []uint
	err     []error
	data    data.Operator
}

func (d Default) User() (uint, []uint) {
	return d.user, d.groups
}

func (m *Default) BindUser() {
	if m.Error() != nil {
		return
	}

	if m.Request == nil {
		log.Println("request not set")
		return
	}

	u := m.Request.Header.Get("X_Session_User")
	user, err := strconv.Atoi(u)
	if err != nil {
		m.Error(err)
		return
	}
	m.user = uint(user)

	if m.Request.Header.Get("X_User_Groups") != "" {
		for _, g := range strings.Split(m.Request.Header.Get("X_User_Groups"), ",") {
			group, err := strconv.Atoi(g)
			if err != nil {
				m.Error(err)
				return
			}
			m.groups = append(m.groups, uint(group))
		}
	}
}

func (m *Default) Bind() {
	if m.Error() != nil {
		return
	}

	if m.Request == nil {
		m.Error("request not set")
		return
	}

	if m.data == nil {
		m.Error("no data set")
		return
	}

	err := m.bind()
	if err != nil {
		m.Error(err)
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
		m.Error(err)
		return
	}
	h := sha1.New()
	h.Write(json)
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Default) Error(es ...interface{}) error {
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

func (m *Default) Manage(db storage.DBManager, action string) {
	if m.Error() != nil {
		return
	}

	db.Manage(m.data, action)
	err := db.Error()
	if err != nil {
		m.Error(err)
		return
	}
}

func (m *Default) Create(db storage.DBCreater) {
	if m.Error() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Error(err)
		return
	}

	db.Create(m.data)
	err = db.Error()
	if err != nil {
		m.Error(err)
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.Error(err)
		return
	}

	dbUpdater := db.CreaterToUpdater()
	dbUpdater.Save(m.data, "with-create")
	err = dbUpdater.Error()
	if err != nil {
		m.Error(err)
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Error(err)
		return
	}

	m.Hasher()
}

func (m *Default) Read(db storage.DBReader, params ...string) {
	if m.Error() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Error(err)
		return
	}

	m.data.Read(db, params...)
	err = db.Error()
	if err != nil {
		m.Error(err)
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Error(err)
		return
	}

	m.Hasher()
}

func (m *Default) Update(db storage.DBUpdater) {
	if m.Error() != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Error(err)
		return
	}

	db.Save(m.data)
	err = db.Error()
	if err != nil {
		m.Error(err)
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Error(err)
		return
	}

	m.Hasher()
}
