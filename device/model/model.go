package model

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
)

type Operator interface {
	Manage(storage.DBManager, string)
	Create(storage.DBCreater)
	Read(storage.DBReader)
	// Update(storage.DBUpdater)
	// Append(string, storage.DBReader)
	Communicator
	DataSelector
}

type Communicator interface {
	Bind()
	User() (uint, []uint)
	Hasher()
	Error() error
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
	err     error
	data    data.Operator
}

func (d Default) User() (uint, []uint) {
	return d.user, d.groups
}

func (m *Default) Bind() {
	if m.err != nil {
		return
	}

	if m.Request == nil {
		m.err = errors.New("request not set")
		return
	}

	u := m.Request.Header.Get("X_Session_User")
	user, err := strconv.Atoi(u)
	if err != nil {
		m.err = err
		return
	}
	m.user = uint(user)

	if m.Request.Header.Get("X_User_Groups") != "" {
		for _, g := range strings.Split(m.Request.Header.Get("X_User_Groups"), ",") {
			group, err := strconv.Atoi(g)
			if err != nil {
				m.err = err
				return
			}
			m.groups = append(m.groups, uint(group))
		}
	}

	if m.data == nil {
		m.err = fmt.Errorf("no data set")
		return
	}

	err = m.bind()
	if err != nil {
		m.err = err
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
		m.err = err
		return
	}
	h := sha1.New()
	h.Write(json)
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Default) Error() error {
	return m.err
}

func (m *Default) Manage(db storage.DBManager, action string) {
	if m.err != nil {
		return
	}

	db.Manage(m.data, action)
	err := db.Error()
	if err != nil {
		m.err = err
		return
	}
}

func (m *Default) Create(db storage.DBCreater) {
	if m.err != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.err = err
		return
	}

	db.Create(m.data)
	err = db.Error()
	if err != nil {
		m.err = err
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.err = err
		return
	}

	dbUpdater := db.CreaterToUpdater()
	dbUpdater.Save(m.data, "with-create")
	err = dbUpdater.Error()
	if err != nil {
		m.err = err
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.err = err
		return
	}

	m.Hasher()
}

func (m *Default) Read(db storage.DBReader) {
	if m.err != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.err = err
		return
	}

	m.data.Read(db)
	err = db.Error()
	if err != nil {
		m.err = err
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.err = err
		return
	}

	m.Hasher()
}
