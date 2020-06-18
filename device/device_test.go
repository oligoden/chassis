package device_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/device"
	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage"
	"github.com/oligoden/chassis/storage/gormdb"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3316)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestMigration(t *testing.T) {
	store := gormdb.New(dbt, uri)

	dMatch := NewDevice(store)
	dMatch.Manage("migrate")

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	if !db.HasTable("users") {
		t.Error(`expected table users`)
	}
	db.Close()
	err = db.Error
	if err != nil {
		t.Error(err)
	}
}

func TestCreate(t *testing.T) {
	cleanDBUserTables()
	setupDBTable("testmodels")

	f := make(url.Values)
	f.Set("field", "test")
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("X_Session_User", `1`)
	w := httptest.NewRecorder()

	s := gormdb.New(dbt, uri)
	d := NewDevice(s)
	d.Manage("migrate")
	d.Create().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	exp := `"ID":1,"Field":"test"`
	got := w.Body.String()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected substring "%s", got "%s"`, exp, got)
	}
}

func TestRead(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables()
	setupDBTable(&TestModel{}, db)

	x := &TestModel{Field: "a"}
	x.Perms = ":::r"
	db.Create(x)
	db.Close()
	err = db.Error
	if err != nil {
		t.Error(err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("X_Session_User", `1`)
	w := httptest.NewRecorder()

	s := gormdb.New(dbt, uri)
	d := NewDevice(s)
	d.Manage("migrate")
	d.Read().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	exp := `"ID":1,"Field":"a"`
	got := w.Body.String()
	if !strings.Contains(got, exp) {
		t.Errorf(`expected substring "%s", got "%s"`, exp, got)
	}
}

type Device struct {
	device.Default
}

func NewDevice(s storage.Storer) *Device {
	d := &Device{}
	nm := func(r *http.Request) model.Operator { return NewModel(r) }
	nv := func(w http.ResponseWriter) view.Operator { return NewView(w) }
	d.Default = device.NewDevice(nm, nv, s)
	return d
}

type Model struct {
	model.Default
}

func NewModel(r *http.Request) *Model {
	m := &Model{}
	m.Default = model.Default{}
	m.Request = r
	m.NewData = func() data.Operator { return NewTestModel() }
	m.Data(NewTestModel())
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

type TestModel struct {
	ID    uint   `gorm:"primary_key"`
	Field string `form:"field"`
	// Players []Player `form:"-" json:"players" gorm:"foreignkey:MatchUC;association_foreignkey:UC"`
	data.Default
}

func NewTestModel() *TestModel {
	r := &TestModel{}
	r.Default = data.Default{}
	r.Perms = "ru:ru:c:"
	r.Groups(2)
	return r
}

func (TestModel) TableName() string {
	return "testmodels"
}

func (x *TestModel) Read(db storage.DBReader) error {
	db.First(x)
	err := db.Error()
	if err != nil {
		return err
	}
	return nil
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
