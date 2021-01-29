package model_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"
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
	exp := `strconv.Atoi: parsing "": invalid syntax`
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
	exp := `strconv.Atoi: parsing "a": invalid syntax`
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
	exp := `strconv.Atoi: parsing "": invalid syntax`
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
	exp := `strconv.Atoi: parsing "a": invalid syntax`
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
	exp := `strconv.Atoi: parsing "a": invalid syntax`
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
	m.NewData = func() data.Operator { return NewTestData() }
	m.Data(NewTestData())
	return m
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
	Field string `form:"field" json:"field"`
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
	q := "CREATE TABLE `testdata` (`field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("doing test_data migration: %w", err)
	}
	return nil
}

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
