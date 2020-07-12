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
	Hasher() error
	TableName() string
	storage.Authenticator
}

type Default struct {
	UC       string `gorm:"unique" json:"uc" form:"uc"`
	GroupIDs []uint `gorm:"-" json:"-"`
	UserIDs  []uint `gorm:"-" json:"-"`
	OwnerID  uint   `json:"-"`
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

func (m *Default) Groups(g ...uint) []uint {
	m.GroupIDs = append(m.GroupIDs, g...)
	return m.GroupIDs
}

func (m *Default) Users(u ...uint) []uint {
	m.UserIDs = append(m.UserIDs, u...)
	return m.UserIDs
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
