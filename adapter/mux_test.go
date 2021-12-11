package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oligoden/chassis/adapter"
	"github.com/steinfletcher/apitest"
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
	mux := adapter.NewMux().
		SetURL("http://test.com:8080").
		Compile(Mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.Host = "test.com:8080"
	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html></html>").
		End()

	req.Host = "staging.test.com:8080"
	apitest.New().
		Handler(mux).
		HttpRequest(req).
		Expect(t).
		Status(http.StatusOK).
		Body("<html>subdomain</html>").
		End()
}
