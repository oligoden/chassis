package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/oligoden/chassis/adapter"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

var Mux func(*adapter.Mux) = func(m *adapter.Mux) {
	// s := m.Stores["mysqldb"]

	// dRouting := routing.NewDevice(s, m.RPDs...)
	// s.Migrate(routing.NewRecord())

	m.Handle("/").
		Core(adapter.ServeFile("static/index.html")).
		SubDomain(adapter.ServeFile("static/subdomain.html"), "-api").
		// And(dSession.Authenticate()).
		Notify().Entry()

	// m.Handle("/profiles").
	// 	NotFound().
	// 	SubDomain(dRouting.Check(), "api").
	// 	And(dSession.CreateUser()).
	// 	And(dSession.Authenticate()).
	// 	CORS().Notify().Entry()
}

func Test(t *testing.T) {
	assert := assert.New(t)
	var err error

	mux := adapter.NewMux().
		SetURL("https://test.com:8080").
		Compile(Mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.URL, err = url.Parse("http://test.com/")
	assert.NoError(err)

	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html></html>").
		End()

	req.URL, err = url.Parse("http://test.com:8080/")
	assert.NoError(err)
	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html></html>").
		End()

	req.URL, err = url.Parse("http://subdomain.test.com/")
	assert.NoError(err)
	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html>subdomain</html>").
		End()
}
