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

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (a Adapter) Notify() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("--> request %s %s from %s\n--  %s\n--  user %s, session %s\n", r.Method, r.URL.String(), r.RemoteAddr, r.UserAgent(), r.Header.Get("X_user"), r.Header.Get("X_session"))

			lw := NewLoggingResponseWriter(w)
			a.Handler.ServeHTTP(lw, r)

			fmt.Printf("<-- response %d %s %s\n\n", lw.statusCode, w.Header().Get("X_user"), w.Header().Get("X_session"))
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

func (a Adapter) Delete(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
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
