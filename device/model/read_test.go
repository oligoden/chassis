package model_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestRead(t *testing.T) {
	testCleanup(t)
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` (`field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X_Session_User", "1")

	s := gosql.New(dbt, uri)
	m := NewModel(req, s)
	e := &TestData{}
	m.Data(e)

	m.Read()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	exp := "a"
	got := e.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
