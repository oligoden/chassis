package view_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/device/view"
)

func TestView(t *testing.T) {
	w := httptest.NewRecorder()

	e := NewTestData()
	e.Field = "a"
	e.Date, _ = time.Parse("2006-01-02", "2021-03-01")

	m := NewModel(nil, nil)
	m.Data(e)
	v := NewView(w)
	v.JSON(m)

	exp := "200"
	got := fmt.Sprint(w.Code)
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = `{"field":"a","uc":"","date":"2021-03-01"}`
	got = w.Body.String()
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

type Model struct {
	model.Default
}

func NewModel(r *http.Request, s model.Connector) *Model {
	m := &Model{}
	m.Default = model.Default{}
	m.Request = r
	m.Store = s
	m.BindUser()
	m.NewData = func() data.Operator { return NewTestData() }
	m.Data(NewTestData())
	return m
}

type View struct {
	view.Default
}

func NewView(w http.ResponseWriter) *View {
	v := &View{}
	v.Default = view.Default{}
	v.Response = w
	return v
}

type TestData struct {
	Field string    `form:"field" json:"field"`
	Date  time.Time `form:"date" json:"date"`
	data.Default
}

func NewTestData() *TestData {
	r := &TestData{}
	r.Default = data.Default{}
	r.Perms = "ru:ru:c:c"
	r.Groups(2)
	return r
}

func (TestData) TableName() string {
	return "testdata"
}

func (e TestData) MarshalJSON() ([]byte, error) {
	type Alias TestData
	return json.Marshal(&struct {
		Alias
		Date string `json:"date"`
	}{
		Alias: (Alias)(e),
		Date:  e.Date.Format("2006-01-02"),
	})
}
