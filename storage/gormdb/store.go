package gormdb

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage"
)

type Store struct {
	dbt              string
	uri              string
	gorm             *gorm.DB
	err              error
	uniqueCodeLength uint
	ucFunc           func(uint, rand.Source) string
	rs               rand.Source
}

func New(dbt, uri string) *Store {
	s := new(Store)

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		s.err = fmt.Errorf("opening db connection for new store migration: %w", err)
		return s
	}
	db.LogMode(true)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Group{})
	// db.AutoMigrate(&UserGroup{})
	db.AutoMigrate(&RecordGroup{})
	db.AutoMigrate(&RecordUser{})
	db.Close()
	err = db.Error
	if err != nil {
		s.err = fmt.Errorf("doing new store db migration: %w", err)
		return s
	}

	s.dbt = dbt
	s.uri = uri

	s.uniqueCodeLength = 2
	s.ucFunc = chassis.RandString
	s.rs = rand.NewSource(time.Now().UnixNano())
	return s
}

func (s *Store) Connect(id storage.Identificator) storage.DoCloser {
	return nil
}

func (m *Store) UniqueCodeLength(ucl ...uint) uint {
	if len(ucl) > 0 {
		m.uniqueCodeLength = ucl[0]
	}
	return m.uniqueCodeLength
}

func (m *Store) UniqueCodeFunc(ucf ...func(uint, rand.Source) string) func(uint, rand.Source) string {
	if len(ucf) > 0 {
		m.ucFunc = ucf[0]
	}
	return m.ucFunc
}

func (s Store) Error() error {
	return s.err
}

type User struct {
	OwnerID    uint      `gorm:"primary_key" json:"-"`
	UC         string    `gorm:"unique" json:"uc" form:"uc"`
	TS         time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Username   string    `gorm:"not null" json:"username"`
	Password   string    `gorm:"-" json:"password"`
	PassHash   string    `json:"-"`
	Salt       string    `json:"salt"`
	UserGroups []Group   `gorm:"many2many:user_groups"`
	GroupIDs   []uint    `gorm:"-" json:"-"`
	UserIDs    []uint    `gorm:"-" json:"-"`
	Perms      string    `json:"-"`
	Hash       string    `json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (e User) Prepare() error {
	return nil
}

func (e *User) Read(db storage.DBReader, params ...string) error {
	db.First(e, params...)
	return nil
}

func (e User) Complete() error {
	return nil
}

func (e *User) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		e.UC = uc[0]
	}
	return e.UC
}

func (e *User) Permissions(p ...string) string {
	if len(p) > 0 {
		e.Perms = p[0]
	}
	return e.Perms
}

func (e *User) Owner(o ...uint) uint {
	if len(o) > 0 {
		e.OwnerID = o[0]
	}
	return e.OwnerID
}

func (e *User) Groups(g ...uint) []uint {
	e.GroupIDs = append(e.GroupIDs, g...)
	return e.GroupIDs
}

func (e *User) Users(u ...uint) []uint {
	e.UserIDs = append(e.UserIDs, u...)
	return e.UserIDs
}

func (e *User) Hasher() error {
	json, err := json.Marshal(e)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(json)
	e.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

type Group struct {
	ID    uint      `gorm:"primary_key"`
	TS    time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Name  string
	Owner uint
	Perms string
}

func (Group) TableName() string {
	return "groups"
}

type RecordGroup struct {
	ID       uint      `gorm:"primary_key"`
	TS       time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	RecordID string
	GroupID  uint
	Owner    uint
	Perms    string
}

type RecordUser struct {
	ID          uint      `gorm:"primary_key"`
	TS          time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Description string
	RecordID    string
	UserID      uint
	Owner       uint
	Perms       string
}
