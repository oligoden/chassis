package gormdb_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage/gormdb"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3316)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestNew(t *testing.T) {
	cleanDBUserTables()
	gormdb.New(dbt, uri)

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	if !db.HasTable("users") {
		t.Error(`expected table users`)
	}
	if !db.HasTable("groups") {
		t.Error(`expected table groups`)
	}
	if !db.HasTable("user_groups") {
		t.Error(`expected table groups`)
	}
	if !db.HasTable("record_groups") {
		t.Error(`expected table groups`)
	}
	db.Close()
	err = db.Error
	if err != nil {
		t.Error(err)
	}
}

func cleanDBUserTables() {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)

	db.DropTableIfExists("users")
	db.DropTableIfExists("groups")
	db.DropTableIfExists("user_groups")
	db.DropTableIfExists("record_groups")
	db.Close()
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}

func setupDBTable(d interface{}, dbs ...*gorm.DB) {
	var db *gorm.DB
	var err error

	if len(dbs) > 0 {
		db = dbs[0]
	}

	if db == nil {
		db, err = gorm.Open(dbt, uri)
		if err != nil {
			log.Fatal(err)
		}
		db.LogMode(true)
		defer db.Close()
	}

	db.DropTableIfExists(d)
	db.AutoMigrate(d)
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}

type TestModel struct {
	TestModelID uint       `gorm:"primary_key"`
	Field       string     `form:"field"`
	SubModels   []SubModel `form:"-" json:"submodels" gorm:"foreignkey:TestModelID;association_foreignkey:TestModelID"`
	UC          string     `gorm:"unique"`
	OwnerID     uint
	groupIDs    []uint
	Perms       string
	Hash        string
}

type TestModels []TestModel

func (TestModel) TableName() string {
	return "testmodels"
}

func (TestModels) TableName() string {
	return "testmodels"
}

func (m *TestModel) UniqueCode(uc ...string) string {
	fmt.Println("code", uc)
	if len(uc) > 0 {
		m.UC = uc[0]
	}
	return m.UC
}

func (m TestModel) Permissions(p ...string) string {
	return m.Perms
}

func (m *TestModel) Owner(o ...uint) uint {
	if len(o) > 0 {
		m.OwnerID = o[0]
	}
	return m.OwnerID
}

func (m *TestModel) Groups(g ...uint) []uint {
	m.groupIDs = append(m.groupIDs, g...)
	return m.groupIDs
}

type SubModel struct {
	SubModelID  uint `gorm:"primary_key"`
	TestModelID uint
	Name        string `form:"name"`
	UC          string `gorm:"unique"`
	OwnerID     uint
	groupIDs    []uint
	Perms       string
	Hash        string
}

func (SubModel) TableName() string {
	return "submodels"
}

func (m *SubModel) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		m.UC = uc[0]
	}
	return m.UC
}

func (m SubModel) Permissions(p ...string) string {
	return m.Perms
}

func (m *SubModel) Owner(o ...uint) uint {
	if len(o) > 0 {
		m.OwnerID = o[0]
	}
	return m.OwnerID
}

func (m *SubModel) Groups(g ...uint) []uint {
	m.groupIDs = append(m.groupIDs, g...)
	return m.groupIDs
}

type BadModel struct {
	ID       uint `gorm:"primary_key"`
	OwnerID  uint
	groupIDs []uint
	Perms    string
}

func (m BadModel) UniqueCode(uc ...string) string {
	return ""
}

func (m BadModel) Permissions(p ...string) string {
	return m.Perms
}

func (m *BadModel) Owner(o ...uint) uint {
	if len(o) > 0 {
		m.OwnerID = o[0]
	}
	return m.OwnerID
}

func (m *BadModel) Groups(g ...uint) []uint {
	m.groupIDs = append(m.groupIDs, g...)
	return m.groupIDs
}
