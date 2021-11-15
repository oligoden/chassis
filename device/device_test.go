package device_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/oligoden/chassis/device"
	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage/gosql"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestCreate(t *testing.T) {
	uri := "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"

	db := testCleanup(t, uri)
	defer db.Close()

	s := gosql.New(uri)
	s.UniqueCodeFunc(func(c uint) string {
		var a string
		for i := uint(0); i < c; i++ {
			a = a + "a"
		}
		return a
	})
	s.Migrate(NewTestData())
	d := NewDevice(s)

	apitest.New().
		Handler(d.Create()).
		Post("/testdata").
		FormData("field", "test").
		FormData("date", "2021-03-01 00:00:00").
		Header("X_user", "1").
		Header("X_session", "1").
		Expect(t).
		Status(http.StatusOK).
		Assert(
			jsonpath.Chain().
				Equal("field", "test").
				Equal("date", "2021-03-01").
				End(),
		).
		End()

	var field, hash string
	err := db.QueryRow("SELECT field,hash from testdata").Scan(&field, &hash)

	assert := assert.New(t)
	if assert.NoError(err) {
		assert.Equal("test", field)
		assert.NotEmpty(hash)
	}
}

func TestRead(t *testing.T) {
	uri := "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"

	db := testCleanup(t, uri)
	defer db.Close()

	s := gosql.New(uri)
	s.Migrate(NewTestData())
	d := NewDevice(s)

	qs := []string{
		"INSERT INTO `testdata` (`field`, `date`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', '2021-03-01', 'xx', 1, ':::', 'xyz')",
	}
	testSetup(db, t, qs...)

	apitest.New().
		Handler(d.Read()).
		Get("/").
		Header("X_user", "1").
		Header("X_session", "1").
		Expect(t).
		Status(http.StatusOK).
		Assert(
			jsonpath.Chain().
				Equal("field", "a").
				Equal("date", "2021-03-01").
				End(),
		).
		End()
}

func TestUpdate(t *testing.T) {
	uri := "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"

	db := testCleanup(t, uri)
	defer db.Close()

	s := gosql.New(uri)
	s.Migrate(NewTestData())
	d := NewDevice(s)

	qs := []string{
		"INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')",
	}
	testSetup(db, t, qs...)

	apitest.New().
		Handler(d.Update()).
		Put("/testdata").
		FormData("field", "b").
		Header("X_user", "1").
		Header("X_session", "1").
		Expect(t).
		Status(http.StatusOK).
		End()

	var field, hash string
	err := db.QueryRow("SELECT field,hash from testdata").Scan(&field, &hash)

	assert := assert.New(t)
	if assert.NoError(err) {
		assert.Equal("b", field)
		assert.NotEmpty(hash)
	}
}

func TestDelete(t *testing.T) {
	uri := "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"

	db := testCleanup(t, uri)
	defer db.Close()

	s := gosql.New(uri)
	s.Migrate(NewTestData())
	d := NewDevice(s)

	qs := []string{
		"INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')",
		"INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'rty')",
	}
	testSetup(db, t, qs...)

	apitest.New().
		Handler(d.Delete()).
		Delete("/testdata/xx").
		Header("X_user", "1").
		Header("X_session", "1").
		Expect(t).
		Status(http.StatusOK).
		End()

	var field, hash string
	err := db.QueryRow("SELECT field,hash from testdata").Scan(&field, &hash)

	assert := assert.New(t)
	if assert.NoError(err) {
		assert.Equal("b", field)
		assert.NotEmpty(hash)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(id) as count from testdata").Scan(&count)

	if assert.NoError(err) {
		assert.Equal(1, count)
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
	Field string    `form:"field" json:"field"`
	Date  time.Time `form:"date" json:"date"`
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

// func (e *TestData) IDValue(id ...uint) uint {
// 	if len(id) > 0 {
// 		e.ID = id[0]
// 	}
// 	return e.ID
// }

func (e TestData) MarshalJSON() ([]byte, error) {
	type Alias TestData
	return json.Marshal(&struct {
		Alias
		Date string `json:"date"`
	}{
		Alias: (Alias)(e),
		Date:  e.Date.Format("2006-01-02"),
	})
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
		_, err := db.Exec(q)
		if err != nil {
			t.Fatal(err)
		}
	}
}
