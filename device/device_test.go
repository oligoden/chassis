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
	cleanDBTables("dataones")

	f := make(url.Values)
	f.Set("field1", "test")
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("ax_session_user", `{"user": 1}`)

	w := httptest.NewRecorder()

	store := gormdb.New(dbt, uri)
	d := NewDevice(store)
	d.Manage("migrate")
	d.Create().ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf(`expected "%d", got "%d"`, http.StatusOK, w.Code)
	}
	if !strings.Contains(w.Body.String(), `"ID":1,"Field1":"test"`) {
		t.Errorf(`expected "%s", got "%s"`, `"ID":1,"Field1":"test"`, w.Body.String())
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
	m.NewData = func() data.Operator { return NewDataOne() }
	m.Data(NewDataOne())
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

type DataOne struct {
	ID     uint   `gorm:"primary_key"`
	Field1 string `form:"field1"`
	// Players []Player `form:"-" json:"players" gorm:"foreignkey:MatchUC;association_foreignkey:UC"`
	data.Default
}

func NewDataOne() *DataOne {
	r := &DataOne{}
	r.Default = data.Default{}
	r.Perms = "ru:ru:c:"
	r.Groups(2)
	return r
}

func (DataOne) TableName() string {
	return "dataones"
}

func (x DataOne) Response() interface{} {
	return x
}

func cleanDBTables(ts ...string) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)

	for _, t := range ts {
		db.DropTableIfExists(t)
	}

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
