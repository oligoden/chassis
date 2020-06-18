package model_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oligoden/chassis/storage/gormdb"
)

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
	m := NewModel(r)
	x = &TestModel{}
	m.Data(x)

	s := gormdb.New(dbt, uri)
	dbRead := s.ReadDB(1, []uint{})

	m.Read(dbRead)
	if m.Error() != nil {
		t.Error(m.Error())
	}

	dbRead.Close()
	if dbRead.Error() != nil {
		t.Error("got error", dbRead.Error())
	}

	exp := "a"
	got := x.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
