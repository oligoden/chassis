package model_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/oligoden/chassis/storage/gosql"
)

// func TestCreateStartError(t *testing.T) {
// 	m := &Model{}
// 	m.Default = model.Default{}

// 	// simulate error, no request set
// 	m.Bind()

// 	s := gormdb.New(dbt, uri)
// 	db := s.CreateDB(1, []uint{})
// 	m.Create(db)

// 	db.Close()
// 	if db.Error() != nil {
// 		t.Error("got error", db.Error())
// 	}

// 	if m.Error() == nil {
// 		t.Error(`expected error`)
// 	}
// }

// func TestCreatePrepareError(t *testing.T) {
// 	m := &Model{}
// 	m.Default = model.Default{}
// 	m.Data(&prepareErrorData{
// 		Default: data.Default{},
// 	})
// 	s := gormdb.New(dbt, uri)
// 	db := s.CreateDB(1, []uint{})
// 	m.Create(db)

// 	db.Close()
// 	if db.Error() != nil {
// 		t.Error("got error", db.Error())
// 	}

// 	if m.Error() == nil {
// 		t.Error(`expected error`)
// 	}
// }

// func TestCreateError(t *testing.T) {
// 	m := &Model{}
// 	m.Default = model.Default{}
// 	m.Data(&createErrorData{
// 		Default: data.Default{},
// 	})
// 	s := gormdb.New(dbt, uri)
// 	db := s.CreateDB(1, []uint{})
// 	m.Create(db)

// 	db.Close()
// 	if db.Error() == nil {
// 		t.Error(`expected error`)
// 	}

// 	if m.Error() == nil {
// 		t.Error(`expected error`)
// 	}
// }

// func TestCreateHashError(t *testing.T) {
// 	setupDBTable(&hashErrorData{})

// 	m := &Model{}
// 	m.Default = model.Default{}
// 	m.Data(&hashErrorData{
// 		Default: data.Default{
// 			Perms: "::c:",
// 		},
// 	})
// 	s := gormdb.New(dbt, uri)
// 	db := s.CreateDB(1, []uint{})
// 	m.Create(db)

// 	db.Close()
// 	if db.Error() != nil {
// 		t.Error("got error", db.Error())
// 	}

// 	if m.Error() == nil {
// 		t.Error(`expected error`)
// 	}
// }

// func TestCreateCompleteError(t *testing.T) {
// 	setupDBTable(&completeErrorData{})

// 	m := &Model{}
// 	m.Default = model.Default{}
// 	m.Data(&completeErrorData{
// 		Default: data.Default{
// 			Perms: "::c:",
// 		},
// 	})
// 	s := gormdb.New(dbt, uri)
// 	db := s.CreateDB(1, []uint{})
// 	m.Create(db)

// 	db.Close()
// 	if db.Error() != nil {
// 		t.Error("got error", db.Error())
// 	}

// 	if m.Error() == nil {
// 		t.Error(`expected error`)
// 	}
// }

func TestCreate(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` ("
	q += " `field` varchar(255),"
	q += " `date` DATETIME NOT NULL DEFAULT '0000-00-00',"
	q += " `id` int unsigned AUTO_INCREMENT,"
	q += " `uc` varchar(255) UNIQUE,"
	q += " `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	f := make(url.Values)
	f.Set("field", "test")
	f.Set("date", "2021-03-01 00:00:00")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/testdatas", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	s := gosql.New(dbt, uri)
	s.UniqueCodeFunc(func(c uint) string {
		var a string
		for i := uint(0); i < c; i++ {
			a = a + "a"
		}
		return a
	})

	m := NewModel(req, s)
	m.Bind()
	m.Create()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	var field, hash string
	var date time.Time
	err = db.QueryRow("SELECT field,date,hash from testdata").Scan(&field, &date, &hash)
	if err != nil {
		t.Error(err)
	}

	exp := "test"
	got := field
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "2021-03-01"
	got = date.Format("2006-01-02")
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "fc1421a39ae43325360fcc9a4677fd5f02ad63b0"
	got = hash
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
