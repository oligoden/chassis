package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestGenInsert(t *testing.T) {
	e := &TestData{
		Field: "test",
	}
	e.ID = 1
	e.UC = "uc"
	e.OwnerID = 1
	e.Perms = ":::"
	e.Hash = "hash"

	c := gosql.NewConnection(1, []uint{})
	c.GenInsert(e)
	q, vs := c.Query()

	exp := "INSERT INTO testdata(field, date, uc, owner_id, perms, hash) VALUES(?, ?, ?, ?, ?, ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[test 1000-01-01 00:00:00 +0000 UTC uc 1 ::: hash]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreate(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` ("
	q += "`field` varchar(255),"
	q += "`date` DATETIME NOT NULL DEFAULT '1000-01-01',"
	q += "`id` int unsigned AUTO_INCREMENT,"
	q += "`uc` varchar(255) UNIQUE,"
	q += "`owner_id` int unsigned,"
	q += "`perms` varchar(255),"
	q += "`hash` varchar(255),"
	q += "PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	s := gosql.New(uri)
	e := &TestData{}
	e.Field = "test"
	e.Perms = ":::c"

	c := s.Connect(1, []uint{})
	c.Create(e)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	var field string
	err = db.QueryRow("SELECT field from testdata").Scan(&field)
	if err != nil {
		t.Error(err)
	}

	exp := "test"
	got := field
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUniqueCodeGeneration(t *testing.T) {
	testCleanup(t)

	s := gosql.New(uri)

	s.UniqueCodeFunc(func(c uint) string {
		var a string
		for i := uint(0); i < c; i++ {
			a = a + "a"
		}
		return a
	})
	s.UniqueCodeLength(1)

	e := &TestData{}
	e.Perms = ":::c"
	s.Migrate(e)
	c := s.Connect(1, []uint{})
	c.Create(e)
	if s.Err() != nil {
		t.Error(s.Err())
	}
	c.Create(e)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	if len(e.UC) <= 1 {
		t.Errorf(`expected "> 1", got "%d"`, len(e.UC))
	}
}
