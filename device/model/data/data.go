package data

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"

	"github.com/oligoden/chassis/storage"
)

type Operator interface {
	Prepare() error
	Complete() error
	Read(storage.DBReader) error
	Response() interface{}
	Hasher() error
	TableName() string
	storage.Authenticator
}

type Default struct {
	UC       string `gorm:"unique" json:"uc" form:"uc"`
	OwnerID  uint   `json:"-"`
	groupIDs []uint
	Perms    string `json:"-"`
	Hash     string `json:"-"`
}

func (m Default) Prepare() error {
	return nil
}

func (m Default) Read(db storage.DBReader) error {
	return nil
}

func (m Default) Complete() error {
	return nil
}

func (m Default) Response() interface{} {
	return m
}

func (m *Default) TableName() string {
	return "models"
}

func (m *Default) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		m.UC = uc[0]
	}
	return m.UC
}

func (m *Default) Permissions(p ...string) string {
	if len(p) > 0 {
		m.Perms = p[0]
	}
	return m.Perms
}

func (m *Default) Owner(o ...uint) uint {
	if len(o) > 0 {
		m.OwnerID = o[0]
	}
	return m.OwnerID
}

func (m *Default) Groups(gs ...uint) []uint {
	m.groupIDs = append(m.groupIDs, gs...)
	return m.groupIDs
}

func (x *Default) Hasher() error {
	json, err := json.Marshal(x)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(json)
	x.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
