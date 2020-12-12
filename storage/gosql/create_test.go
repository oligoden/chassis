package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestGenInsert(t *testing.T) {
	e := &TestData{
		ID:    1,
		Field: "test",
	}
	e.UC = "uc"
	e.OwnerID = 1
	e.Perms = ":::"
	e.Hash = "hash"

	c := gosql.NewConnection(1, []uint{})
	c.GenInsert(e)
	q, vs := c.Query()

	exp := "INSERT INTO testdata(field, uc, owner_id, perms, hash) VALUES(?, ?, ?, ?, ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[test uc 1 ::: hash]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreate(t *testing.T) {
	testCleanup(t)

	s := gosql.New(dbt, uri)
	e := &TestData{}
	e.Field = "test"
	e.Perms = ":::c"
	s.Migrate(e)
	c := s.Connect(1, []uint{})
	c.Create(e)
	if s.Err() != nil {
		t.Error(s.Err())
	}

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

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

	s := gosql.New(dbt, uri)

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
	if s.Err() != nil {
		t.Error(s.Err())
	}

	if len(e.UC) <= 1 {
		t.Errorf(`expected "> 1", got "%d"`, len(e.UC))
	}
}
