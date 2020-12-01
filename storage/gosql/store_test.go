package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func testCleanup(t *testing.T) {
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE groups")
	db.Exec("DROP TABLE record_groups")
	db.Exec("DROP TABLE record_users")

	db.Exec("DROP TABLE testdata")
}

func TestNewStore(t *testing.T) {
	testCleanup(t)

	s := gosql.New(dbt, uri)

	usersExist := false
	groupsExist := false
	recordGroupsExist := false
	recordUsersExist := false

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	rows, err := db.Query("SHOW tables")
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Error(err)
		}

		if name == "users" {
			usersExist = true
		}
		if name == "groups" {
			groupsExist = true
		}
		if name == "record_groups" {
			recordGroupsExist = true
		}
		if name == "record_users" {
			recordUsersExist = true
		}
	}
	if err := rows.Err(); err != nil {
		t.Error(err)
	}

	if !usersExist {
		t.Error("users table not found")
	}
	if !groupsExist {
		t.Error("groups table not found")
	}
	if !recordGroupsExist {
		t.Error("record_groups table not found")
	}
	if !recordUsersExist {
		t.Error("record_users table not found")
	}

	exp := "2"
	got := fmt.Sprint(s.UniqueCodeLength())
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	s.UniqueCodeLength(3)
	exp = "3"
	got = fmt.Sprint(s.UniqueCodeLength())
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	if len(s.UniqueCodeFunc()(4)) != 4 {
		t.Error("random code length not 4")
	}

	s.UniqueCodeFunc(func(c uint) string {
		return "aaa"
	})

	exp = "aaa"
	got = s.UniqueCodeFunc()(4)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

type TestData struct {
	ID            uint   `gosql:"primary_key"`
	Field         string `form:"field"`
	field         string
	SubData       []SubData `gosql:"-"`
	Many2ManyData []SubData `gosql:"-"`
	UC            string    `gosql:"unique"`
	GroupIDs      []uint    `gosql:"-"`
	UserIDs       []uint    `gosql:"-"`
	OwnerID       uint
	Perms         string
	Hash          string
}

func (TestData) TableName() string {
	return "testdata"
}

func (TestData) Migrate(db *sql.DB) error {
	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("doing test_data migration: %w", err)
	}
	return nil
}

func (e TestData) Permissions(p ...string) string {
	return e.Perms
}

func (e *TestData) Owner(o ...uint) uint {
	if len(o) > 0 {
		e.OwnerID = o[0]
	}
	return e.OwnerID
}

func (e *TestData) Users(u ...uint) []uint {
	e.UserIDs = append(e.UserIDs, u...)
	return e.UserIDs
}

func (e *TestData) Groups(g ...uint) []uint {
	e.GroupIDs = append(e.GroupIDs, g...)
	return e.GroupIDs
}

func (e *TestData) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		e.UC = uc[0]
	}
	return e.UC
}

type TestDataMap map[string]TestData

func (TestDataMap) TableName() string {
	return "testdata"
}

func (e TestDataMap) Permissions(p ...string) string {
	return ""
}

func (e TestDataMap) Owner(o ...uint) uint {
	return 0
}

func (e TestDataMap) Users(u ...uint) []uint {
	return []uint{}
}

func (e TestDataMap) Groups(g ...uint) []uint {
	return []uint{}
}

func (e TestDataMap) UniqueCode(uc ...string) string {
	return ""
}

type SubData struct {
	SubDataID  uint `gorm:"primary_key"`
	TestDataID uint
	Field      string `form:"field"`
	UC         string `gorm:"unique"`
	GroupIDs   []uint `gorm:"-" json:"-"`
	UserIDs    []uint `gorm:"-" json:"-"`
	OwnerID    uint
	Perms      string
	Hash       string
}

// func (SubData) TableName() string {
// 	return "subdata"
// }

// func (m *SubData) UniqueCode(uc ...string) string {
// 	if len(uc) > 0 {
// 		m.UC = uc[0]
// 	}
// 	return m.UC
// }

// func (m SubData) Permissions(p ...string) string {
// 	return m.Perms
// }

// func (m *SubData) Owner(o ...uint) uint {
// 	if len(o) > 0 {
// 		m.OwnerID = o[0]
// 	}
// 	return m.OwnerID
// }

// func (m *SubData) Groups(g ...uint) []uint {
// 	m.GroupIDs = append(m.GroupIDs, g...)
// 	return m.GroupIDs
// }

// func (m *SubData) Users(u ...uint) []uint {
// 	m.UserIDs = append(m.UserIDs, u...)
// 	return m.UserIDs
// }
