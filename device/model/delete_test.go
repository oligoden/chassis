package model_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestDelete(t *testing.T) {
	uri := "chassis:password@tcp(localhost:3309)/chassis?charset=utf8&parseTime=True&loc=Local"

	db := testCleanup(t, uri)
	defer db.Close()

	qs := []string{}

	q := "CREATE TABLE `testdata` ("
	q += " `field` varchar(255),"
	q += " `date` DATETIME NOT NULL DEFAULT '0000-00-00',"
	q += " `id` int unsigned AUTO_INCREMENT,"
	q += " `uc` varchar(255) UNIQUE,"
	q += " `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	qs = append(qs, q)

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	qs = append(qs, q)

	testSetup(db, t, qs...)

	f := make(url.Values)
	f.Set("field", "test")
	req := httptest.NewRequest(http.MethodDelete, "/testdatas/xx", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	s := gosql.New(uri)
	m := NewModel(req, s)
	e := &TestData{}
	m.Data(e)

	m.Bind()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	exp := "xx"
	got := e.UC
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m.Read()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	exp = "a"
	got = e.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m.Delete()
	if m.Err() != nil {
		t.Error(m.Err())
	}

	var field string
	err := db.QueryRow("SELECT field from testdata").Scan(&field)
	if err == nil {
		t.Error("expected no results")
	}
}
