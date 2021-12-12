package adapter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/oligoden/chassis/adapter"
	"github.com/steinfletcher/apitest"
)

var Mux func(*adapter.Mux) = func(m *adapter.Mux) {
	m.Handle("/").
		Core(adapter.ServeFile("static/index.html")).
		SubDomain(adapter.ServeFile("static/subdomain.html"), "-api").
		CORS().Notify().Entry()
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
		HeaderPresent("Access-Control-Allow-Origin").
		HeaderPresent("Access-Control-Allow-Credentials").
		Body("<html>subdomain</html>").
		End()
}
