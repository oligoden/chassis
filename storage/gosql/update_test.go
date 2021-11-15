package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestGenUpdate(t *testing.T) {
	e := &TestData{
		Field: "b",
	}
	e.ID = 1
	e.Perms = ":::"
	e.OwnerID = 1
	e.Hash = "abc"

	c := gosql.NewConnection(1, []uint{})
	c.GenUpdate(e)
	q, vs := c.Query()

	exp := "UPDATE testdata SET field = ?, date = ?, hash = ? WHERE id = ?"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[b 1000-01-01 00:00:00 +0000 UTC abc 1]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUpdate(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
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

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jkl')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	e := &TestData{}

	s := gosql.New(uri)
	c := s.Connect(1, []uint{})
	w := gosql.NewWhere("uc = ?", "yy")
	c.AddModifiers(w)
	c.Read(e)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	e.Field = "c"
	c.Update(e)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	var field string
	err = db.QueryRow("SELECT field FROM testdata WHERE uc = 'yy'").Scan(&field)
	if err != nil {
		t.Error(err)
	}

	exp := "c"
	got := field
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
