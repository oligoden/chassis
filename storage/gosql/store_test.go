package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local&parseTime=true"
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
	db.Exec("DROP TABLE subdata")
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
	Field         string `form:"field"`
	field         string
	SubData       []SubData `gosql:"-"`
	Many2ManyData []SubData `gosql:"-"`
	Default
}

type Default struct {
	ID       uint   `gosql:"primary_key"`
	UC       string `json:"uc" form:"uc"`
	GroupIDs []uint `gosql:"-" json:"-"`
	UserIDs  []uint `gosql:"-" json:"-"`
	OwnerID  uint   `json:"-"`
	Perms    string `json:"-"`
	Hash     string `json:"-"`
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

func (e *TestData) IDValue(id ...uint) uint {
	if len(id) > 0 {
		e.ID = id[0]
	}
	return e.ID
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

func (TestDataMap) IDValue(...uint) uint {
	return 0
}

func (e TestDataMap) UniqueCode(uc ...string) string {
	return ""
}

type TestDataSlice []TestData

func (TestDataSlice) TableName() string {
	return "testdata"
}

func (e TestDataSlice) Permissions(p ...string) string {
	return ""
}

func (e TestDataSlice) Owner(o ...uint) uint {
	return 0
}

func (e TestDataSlice) Users(u ...uint) []uint {
	return []uint{}
}

func (e TestDataSlice) Groups(g ...uint) []uint {
	return []uint{}
}

func (TestDataSlice) IDValue(...uint) uint {
	return 0
}

func (e TestDataSlice) UniqueCode(uc ...string) string {
	return ""
}

type SubData struct {
	TestDataID uint
	Field      string `form:"field"`
	Default
}

func (SubData) TableName() string {
	return "subdata"
}

func (e *SubData) IDValue(id ...uint) uint {
	if len(id) > 0 {
		e.ID = id[0]
	}
	return e.ID
}

func (m *SubData) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		m.UC = uc[0]
	}
	return m.UC
}

func (m SubData) Permissions(p ...string) string {
	return m.Perms
}

func (m *SubData) Owner(o ...uint) uint {
	if len(o) > 0 {
		m.OwnerID = o[0]
	}
	return m.OwnerID
}

func (m *SubData) Groups(g ...uint) []uint {
	m.GroupIDs = append(m.GroupIDs, g...)
	return m.GroupIDs
}

func (m *SubData) Users(u ...uint) []uint {
	m.UserIDs = append(m.UserIDs, u...)
	return m.UserIDs
}

type SubDataMap map[string]SubData

func (SubDataMap) TableName() string {
	return "subdata"
}

func (e SubDataMap) Permissions(p ...string) string {
	return ""
}

func (e SubDataMap) Owner(o ...uint) uint {
	return 0
}

func (e SubDataMap) Users(u ...uint) []uint {
	return []uint{}
}

func (e SubDataMap) Groups(g ...uint) []uint {
	return []uint{}
}

func (SubDataMap) IDValue(...uint) uint {
	return 0
}

func (e SubDataMap) UniqueCode(uc ...string) string {
	return ""
}

type SubDataMapID map[uint]SubData

func (SubDataMapID) TableName() string {
	return "subdata"
}

func (e SubDataMapID) Permissions(p ...string) string {
	return ""
}

func (e SubDataMapID) Owner(o ...uint) uint {
	return 0
}

func (e SubDataMapID) Users(u ...uint) []uint {
	return []uint{}
}

func (e SubDataMapID) Groups(g ...uint) []uint {
	return []uint{}
}

func (SubDataMapID) IDValue(...uint) uint {
	return 0
}

func (e SubDataMapID) UniqueCode(uc ...string) string {
	return ""
}
