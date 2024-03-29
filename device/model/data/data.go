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
	Hasher() error
	storage.Operator
}

type Default struct {
	ID       uint   `gosql:"primary_key" json:"-"`
	UC       string `gorm:"unique" json:"uc" form:"uc"`
	GroupIDs []uint `gosql:"-" json:"-"`
	UserIDs  []uint `gosql:"-" json:"-"`
	OwnerID  uint   `json:"-"`
	Perms    string `json:"-"`
	Hash     string `json:"-"`
}

func (e Default) Prepare() error {
	return nil
}

func (e Default) Complete() error {
	return nil
}

func (e *Default) TableName() string {
	return "models"
}

func (e *Default) IDValue(id ...uint) uint {
	if len(id) > 0 {
		e.ID = id[0]
	}
	return e.ID
}

func (e *Default) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		e.UC = uc[0]
	}
	return e.UC
}

func (e *Default) Permissions(p ...string) string {
	if len(p) > 0 {
		e.Perms = p[0]
	}
	return e.Perms
}

func (e *Default) Owner(o ...uint) uint {
	if len(o) > 0 {
		e.OwnerID = o[0]
	}
	return e.OwnerID
}

func (e *Default) Groups(g ...uint) []uint {
	e.GroupIDs = append(e.GroupIDs, g...)
	return e.GroupIDs
}

func (e *Default) Users(u ...uint) []uint {
	e.UserIDs = append(e.UserIDs, u...)
	return e.UserIDs
}

func (e *Default) Hasher() error {
	json, err := json.Marshal(e)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(json)
	e.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}
