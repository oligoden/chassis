package adapter

import (
	"fmt"
	"net/http"
)

type Adapter struct {
	Handler http.Handler
}

func (a Adapter) Entry() http.Handler {
	return a.Handler
}

func (a Adapter) Notify() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("--> request %s %s from %s\n\t%s\n\t%s\n", r.Method, r.URL.String(), r.RemoteAddr, r.UserAgent(), r.Header.Get("ax_session_user"))
			a.Handler.ServeHTTP(w, r)
			fmt.Printf("<-- response %s\n\n", w.Header().Get("X_Session_User"))
		}),
	}
}

func (a Adapter) And(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			a.Handler.ServeHTTP(w, r)
		}),
	}
}

func (a Adapter) Get(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
	}
}

func (a Adapter) Post(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
	}
}

func (a Adapter) Put(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
	}
}

func Core(h http.Handler) Adapter {
	return Adapter{
		Handler: h,
	}
}

func MNA() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}),
	}
}
