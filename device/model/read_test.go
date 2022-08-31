package model_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestRead(t *testing.T) {
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

	q = "INSERT INTO `testdata` (`field`, `date`, `uc`, `owner_id`, `perms`, `hash`)"
	q += " VALUES ('a', '2021-03-01', 'xx', 1, ':::', 'xyz')"
	qs = append(qs, q)

	q = "INSERT INTO `testdata` (`field`, `date`, `uc`, `owner_id`, `perms`, `hash`)"
	q += " VALUES ('b', '2021-03-01', 'yy', 1, ':::', 'dfg')"
	qs = append(qs, q)

	testSetup(db, t, qs...)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	s := gosql.New(uri)
	m := NewModel(req, s)
	e := &TestData{}
	m.Data(e)

	m.Read()

	assert := assert.New(t)
	assert.NoError(m.Err())
	assert.NotEmpty(e.Field)
	assert.Equal("2021-03-01", e.Date.Format("2006-01-02"))

	req = httptest.NewRequest(http.MethodGet, "/testdata/xx", nil)
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	m = NewModel(req, s)
	e = &TestData{}
	m.Data(e)
	m.Bind()
	m.Read()

	assert.NoError(m.Err())
	assert.Equal("a", e.Field)

	req = httptest.NewRequest(http.MethodGet, "/testdata/yy", nil)
	req.Header.Set("X_user", "1")
	req.Header.Set("X_session", "1")

	m = NewModel(req, s)
	e = &TestData{}
	m.Data(e)
	m.Bind()
	m.Read()

	assert.NoError(m.Err())
	assert.Equal("b", e.Field)
}
