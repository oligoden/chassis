package device_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oligoden/chassis/device"
	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage/gosql"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestCreate(t *testing.T) {
	testCleanup(t)

	f := make(url.Values)
	f.Set("field", "test")
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("X_user", `1`)
	r.Header.Set("X_session", `1`)
	w := httptest.NewRecorder()

	s := gosql.New(dbt, uri)
	s.UniqueCodeFunc(func(c uint) string {
		var a string
		for i := uint(0); i < c; i++ {
			a = a + "a"
		}
		return a
	})
	s.Migrate(NewTestData())
	d := NewDevice(s)
	d.Create().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	exp := `"field":"test"`
	got := w.Body.String()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected substring "%s", got "%s"`, exp, got)
	}

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	var field, hash string
	err = db.QueryRow("SELECT field,hash from testdata").Scan(&field, &hash)
	if err != nil {
		t.Error(err)
	}

	exp = "test"
	got = field
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "fc1421a39ae43325360fcc9a4677fd5f02ad63b0"
	got = hash
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestRead(t *testing.T) {
	testCleanup(t)
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` (`field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X_user", `1`)
	r.Header.Set("X_session", `1`)
	w := httptest.NewRecorder()

	s := gosql.New(dbt, uri)
	s.Migrate(NewTestData())
	d := NewDevice(s)
	d.Read().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	exp := `"field":"a"`
	got := w.Body.String()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected substring "%s", got "%s"`, exp, got)
	}
}

func TestUpdate(t *testing.T) {
	testCleanup(t)
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` (`field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	f := make(url.Values)
	f.Set("field", "b")
	r := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(f.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("X_user", `1`)
	r.Header.Set("X_session", `1`)
	w := httptest.NewRecorder()

	s := gosql.New(dbt, uri)
	s.Migrate(NewTestData())
	d := NewDevice(s)
	d.Update().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	exp := `field":"b"`
	got := w.Body.String()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected substring "%s", got "%s"`, exp, got)
	}
}

type Device struct {
	device.Default
}

func NewDevice(s model.Connector) *Device {
	d := &Device{}

	nm := func(r *http.Request) model.Operator {
		return NewModel(r, s)
	}

	nv := func(w http.ResponseWriter) view.Operator {
		return NewView(w)
	}

	d.Default = device.NewDevice(nm, nv, s)
	return d
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
	r.Perms = "ru:ru:c:"
	r.Groups(2)
	return r
}

func (TestData) TableName() string {
	return "testdata"
}

func (e *TestData) IDValue(id ...uint) uint {
	if len(id) > 0 {
		e.ID = id[0]
	}
	return e.ID
}

func (TestData) Migrate(db *sql.DB) error {
	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
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
