package model_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
)

func TestDataSetting(t *testing.T) {
	e := &TestData{}
	m := &Model{}
	m.Default = model.Default{}

	if m.Data() != nil {
		t.Errorf(`expected nil`)
	}

	if m.Data(e) == nil {
		t.Errorf(`expected not nil`)
	}
}

func TestHashing(t *testing.T) {
	e := &TestData{}
	m := &Model{}
	m.Default = model.Default{}
	m.Data(e)

	if m.Hash != "" {
		t.Errorf(`expected empty hash`)
	}

	m.Hasher()

	if m.Hash == "" {
		t.Errorf(`expected non-empty hash`)
	}
}

func TestBindStartError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}

	m.Bind()
	exp := "request not set"
	got := m.Err().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	// calling Bind() with existing error should return immediately
	m.Bind()
	exp = "request not set"
	got = m.Err().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindUserNoUserError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{
		Request: httptest.NewRequest(http.MethodPost, "/", nil),
	}

	m.BindUser()
	exp := `X_user not set`
	got := m.Err().Error()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindUserNotUserIntError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_user", "a")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `user binding X_user`
	got := m.Err().Error()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindUserNoSessionError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_user", "1")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `X_session not set`
	got := m.Err().Error()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindNotSessionIntError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "a")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `session binding X_session`
	got := m.Err().Error()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindNotGroupIntError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")
	req.Header.Set("X_User_Groups", "a")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `user binding X_user_groups`
	got := m.Err().Error()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindNoDataError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	m.Default = model.Default{
		Request: req,
	}

	m.Bind()
	exp := "no data set"
	got := m.Err().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

type Model struct {
	model.Default
}

func NewModel(r *http.Request, s model.Connector) *Model {
	m := &Model{}
	m.Default = model.Default{}
	m.Request = r
	m.Store = s
	m.BindUser()
	m.Data(NewTestData())
	return m
}

func (m *Model) NewData(ds ...string) {
	if len(ds) > 0 {
		if ds[0] == "list" {
			m.Data(NewTestDataList())
		}
	}
}

type View struct {
	view.Default
}

func NewView(w http.ResponseWriter) *View {
	v := &View{}
	v.Default = view.Default{}
	v.Response = w
	return v
}

type TestData struct {
	Field string    `form:"field" json:"field"`
	Date  time.Time `form:"date" json:"date"`
	data.Default
}

func NewTestData() *TestData {
	r := &TestData{}
	r.Default = data.Default{}
	r.Perms = "ru:ru:c:c"
	r.Groups(2)
	return r
}

func (TestData) TableName() string {
	return "testdata"
}

func (TestData) Migrate(db *sql.DB) error {
	q := "CREATE TABLE `testdata` ("
	q += " `field` varchar(255),"
	q += " `date` DATETIME NOT NULL DEFAULT '0000-00-00',"
	q += " `id` int unsigned AUTO_INCREMENT,"
	q += " `uc` varchar(255) UNIQUE,"
	q += " `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("doing test_data migration: %w", err)
	}
	return nil
}

type TestDataList map[string]TestData

func NewTestDataList() TestDataList {
	return TestDataList{}
}

func (TestDataList) TableName() string {
	return "testdata"
}

func (TestDataList) Complete() error {
	return nil
}

func (TestDataList) Hasher() error {
	return nil
}

func (TestDataList) Prepare() error {
	return nil
}

func (e TestDataList) Users(u ...uint) []uint {
	return []uint{}
}

func (e TestDataList) Groups(g ...uint) []uint {
	return []uint{}
}

func (TestDataList) IDValue(...uint) uint {
	return 0
}

func (e TestDataList) Owner(o ...uint) uint {
	return 0
}

func (e TestDataList) Permissions(p ...string) string {
	return ""
}

func (e TestDataList) UniqueCode(uc ...string) string {
	return ""
}

func testCleanup(t *testing.T, uri string) *sql.DB {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		t.Error(err)
	}

	db.Exec("DROP TABLE users")
	db.Exec("DROP TABLE groups")
	db.Exec("DROP TABLE record_groups")
	db.Exec("DROP TABLE record_users")

	db.Exec("DROP TABLE testdata")

	return db
}

func testSetup(db *sql.DB, t *testing.T, qs ...string) {
	for _, q := range qs {
		fmt.Println("running", q)
		_, err := db.Exec(q)
		if err != nil {
			t.Fatal(err)
		}
	}
}
