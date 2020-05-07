package gormdb

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/device/model/data"
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
	ID       uint      `gorm:"primary_key" form:"id"`
	TS       time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Username string    `gorm:"unique;not null" json:"username"`
	Password string    `gorm:"-" json:"password"`
	Salt     string    `json:"salt"`
	Groups   []Group   `gorm:"many2many:user_groups;foreignkey:id;association_foreignkey:id"`
	*data.Default
}

func (User) TableName() string {
	return "users"
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
	GroupID  uint
	RecordID string
	Owner    uint
	Perms    string
}
