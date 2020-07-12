package model_test

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3316)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestData(t *testing.T) {
	xTestModel := &TestModel{}
	mTestModel := &Model{}
	mTestModel.Default = model.Default{}

	if mTestModel.Data() != nil {
		t.Errorf(`expected nil`)
	}

	if mTestModel.Data(xTestModel) == nil {
		t.Errorf(`expected not nil`)
	}
}

func TestHashing(t *testing.T) {
	xTestModel := &TestModel{}
	mTestModel := &Model{}
	mTestModel.Default = model.Default{}
	mTestModel.Data(xTestModel)

	if mTestModel.Hash != "" {
		t.Errorf(`expected empty hash`)
	}

	mTestModel.Hasher()

	if mTestModel.Hash == "" {
		t.Errorf(`expected non-empty hash`)
	}
}

func TestBindStartError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}

	m.Bind()
	exp := "request not set"
	got := m.Error().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	// calling Bind() with existing error should return immediately
	m.Bind()
	exp = "request not set"
	got = m.Error().Error()
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
	got := m.Error().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindUserNotUserIntError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_Session_User", "a")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `strconv.Atoi: parsing "a": invalid syntax`
	got := m.Error().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestBindNotGroupIntError(t *testing.T) {
	m := &Model{}
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X_Session_User", "1")
	req.Header.Set("X_User_Groups", "a")
	m.Default = model.Default{
		Request: req,
	}

	m.BindUser()
	exp := `strconv.Atoi: parsing "a": invalid syntax`
	got := m.Error().Error()
	if got != exp {
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
	got := m.Error().Error()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
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

type prepareErrorData struct {
	data.Default
}

func (m prepareErrorData) Prepare() error {
	return errors.New("error")
}

type createErrorData struct {
	data.Default
}

type hashErrorData struct {
	data.Default
}

func (m hashErrorData) Hasher() error {
	return errors.New("hash test error")
}

type completeErrorData struct {
	data.Default
}

func (m completeErrorData) Complete() error {
	return errors.New("complete test error")
}

func NewTestModel() data.Operator {
	r := &TestModel{}
	r.Default = data.Default{}
	r.Permissions("::c:")
	return r
}

type TestModel struct {
	ID    uint   `gorm:"primary_key"`
	Field string `form:"field"`
	// CaptainUC      string  `json:"-"`
	// Captain        Player            `json:"kaptain" gorm:"foreignkey:KaptainUC;association_foreignkey:UC"`
	// Players        []Player          `json:"-" gorm:"foreignkey:MatchUC;association_foreignkey:UC"`
	// PlayerMap      map[string]Player `json:"players" gorm:"-"`
	data.Default
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
