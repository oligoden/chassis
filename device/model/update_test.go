package model_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestUpdate(t *testing.T) {
	testCleanup(t)
	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
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

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	f := make(url.Values)
	f.Set("field", "test")
	req := httptest.NewRequest(http.MethodPut, "/testdatas", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

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

	m.Bind()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	m.Update()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	var field, hash string
	err = db.QueryRow("SELECT field,hash from testdata").Scan(&field, &hash)
	if err != nil {
		t.Error(err)
	}

	exp = "test"
	got = field
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "8f0a824fa3a483940710071db416dab40a16a6ed"
	got = hash
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
