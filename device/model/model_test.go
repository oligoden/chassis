package model_test

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage/gormdb"
)

const (
	dbt = "mysql"
	uri = "chassis:password@tcp(localhost:3316)/chassis?charset=utf8&parseTime=True&loc=Local"
)

func TestData(t *testing.T) {
	xMatch := &Match{}
	mMatch := &Model{}
	mMatch.Default = model.Default{}

	if mMatch.Data() != nil {
		t.Errorf(`expected nil`)
	}

	if mMatch.Data(xMatch) == nil {
		t.Errorf(`expected not nil`)
	}
}

func TestHashing(t *testing.T) {
	xMatch := &Match{}
	mMatch := &Model{}
	mMatch.Default = model.Default{}
	mMatch.Data(xMatch)

	if mMatch.Hash != "" {
		t.Errorf(`expected empty hash`)
	}

	mMatch.Hasher()

	if mMatch.Hash == "" {
		t.Errorf(`expected non-empty hash`)
	}
}

func TestBindStartError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{Err: errors.New("error")}
	m.Bind()
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestBindNoDataError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}
	m.Bind()
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreateStartError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{Err: errors.New("error")}
	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreatePrepareError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}
	m.Data(&prepareErrorData{
		Default: data.Default{},
	})
	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreateError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}
	m.Data(&createErrorData{
		Default: data.Default{},
	})
	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreateHashError(t *testing.T) {
	setupDBTable(&hashErrorData{})

	m := &Model{}
	m.Default = model.Default{}
	m.Data(&hashErrorData{
		Default: data.Default{
			Perms: "::c:",
		},
	})
	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreateCompleteError(t *testing.T) {
	setupDBTable(&completeErrorData{})

	m := &Model{}
	m.Default = model.Default{}
	m.Data(&completeErrorData{
		Default: data.Default{
			Perms: "::c:",
		},
	})
	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)
	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreate(t *testing.T) {
	setupDBTable(&Match{})

	f := make(url.Values)
	f.Set("field", "Chesterfield")
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	m := &Model{}
	m.Default = model.Default{}
	m.Request = req
	m.Data(NewMatch())

	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})

	m.Bind()
	m.Create(db)
	if m.Error() != nil {
		t.Error(m.Error())
	}

	db.Close()
	if db.Error() != nil {
		t.Error("got error", db.Error())
	}

	dbGorm, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}

	xMatch := &Match{}
	dbGorm.First(xMatch)

	exp := "Chesterfield"
	got := xMatch.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
	if xMatch.UC == "" {
		t.Error(`expected non empty unique code`)
	}
	if xMatch.Hash == "" {
		t.Error(`expected non empty hash`)
	}
}

type Model struct {
	model.Default
}

func NewModel(w http.ResponseWriter, r *http.Request) *Model {
	m := &Model{}
	m.Default = model.Default{}
	m.Request = r
	m.NewData = func() data.Operator { return NewMatch() }
	m.Data(NewMatch())
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

func NewMatch() data.Operator {
	r := &Match{}
	r.Default = data.Default{}
	r.Permissions("::c:")
	return r
}

type Match struct {
	ID             uint    `gorm:"primary_key"`
	Field          string  `form:"field"`
	AvgTemperature float32 `form:"avgtemp"`
	Capacity       uint    `form:"capacity"`
	OpenRoof       bool    `form:"openroof"`
	// CaptainUC      string  `json:"-"`
	// Captain        Player            `json:"kaptain" gorm:"foreignkey:KaptainUC;association_foreignkey:UC"`
	// Players        []Player          `json:"-" gorm:"foreignkey:MatchUC;association_foreignkey:UC"`
	// PlayerMap      map[string]Player `json:"players" gorm:"-"`
	data.Default
}

func (Match) TableName() string {
	return "matches"
}

func setupDBTable(d interface{}) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)

	db.DropTableIfExists(d)
	db.AutoMigrate(d)
	db.Close()
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}
