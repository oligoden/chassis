package model_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestBindTime(t *testing.T) {
	f := make(url.Values)
	f.Set("date", "2021-03-01 00:00:00")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tests", strings.NewReader(f.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	m := NewModel(req, nil)
	e := NewTestData()
	m.Data(e)
	m.Bind()
	if m.Err() != nil {
		t.Fatal(m.Err())
	}

	exp := "2021-03-01"
	got := e.Date.Format("2006-01-02")
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
