package model

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

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
	Err     error                `json:"-"`
	data    data.Operator
}

func (m *Default) Bind() {
	if m.Err != nil {
		return
	}

	if m.data == nil {
		m.Err = fmt.Errorf("no data set")
		return
	}

	err := m.bind()
	if err != nil {
		m.Err = err
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
		m.Err = err
		return
	}
	h := sha1.New()
	h.Write(json)
	m.Hash = fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Default) Error() error {
	return m.Err
}

func (m *Default) Manage(db storage.DBManager, action string) {
	if m.Err != nil {
		return
	}

	db.Manage(m.data, action)
	err := db.Error()
	if err != nil {
		m.Err = err
		return
	}
}

func (m *Default) Create(db storage.DBCreater) {
	if m.Err != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Err = err
		return
	}

	db.Create(m.data)
	err = db.Error()
	if err != nil {
		m.Err = err
		return
	}

	err = m.data.Hasher()
	if err != nil {
		m.Err = err
		return
	}

	dbUpdater := db.CreaterToUpdater()
	dbUpdater.Save(m.data, "with-create")
	err = dbUpdater.Error()
	if err != nil {
		m.Err = err
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Err = err
		return
	}

	m.Hasher()
}

func (m *Default) Read(db storage.DBReader) {
	if m.Err != nil {
		return
	}

	err := m.data.Prepare()
	if err != nil {
		m.Err = err
		return
	}

	m.data.Read(db)
	err = db.Error()
	if err != nil {
		m.Err = err
		return
	}

	err = m.data.Complete()
	if err != nil {
		m.Err = err
		return
	}

	m.Hasher()
}
