package model_test

import (
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

func TestCreateStartError(t *testing.T) {
	m := &Model{}
	m.Default = model.Default{}

	// simulate error, no request set
	m.Bind()

	s := gormdb.New(dbt, uri)
	db := s.CreateDB(1, []uint{})
	m.Create(db)

	db.Close()
	if db.Error() != nil {
		t.Error("got error", db.Error())
	}

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

	db.Close()
	if db.Error() != nil {
		t.Error("got error", db.Error())
	}

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

	db.Close()
	if db.Error() == nil {
		t.Error(`expected error`)
	}

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

	db.Close()
	if db.Error() != nil {
		t.Error("got error", db.Error())
	}

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

	db.Close()
	if db.Error() != nil {
		t.Error("got error", db.Error())
	}

	if m.Error() == nil {
		t.Error(`expected error`)
	}
}

func TestCreate(t *testing.T) {
	dbGorm, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	dbGorm.LogMode(true)
	setupDBTable(&TestModel{}, dbGorm)

	f := make(url.Values)
	f.Set("field", "test")
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X_Session_User", "1")

	m := &Model{}
	m.Default = model.Default{}
	m.Request = req
	m.Data(NewTestModel())

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

	xTestModel := &TestModel{}
	dbGorm.First(xTestModel)
	dbGorm.Close()
	err = dbGorm.Error
	if err != nil {
		t.Error(err)
	}

	exp := "test"
	got := xTestModel.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
	if xTestModel.UC == "" {
		t.Error(`expected non empty unique code`)
	}
	if xTestModel.Hash == "" {
		t.Error(`expected non empty hash`)
	}
}
