package adapter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"xojoc.pw/useragent"
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
			params := []interface{}{r.Method, r.URL.String(), r.RemoteAddr}
			text := "--> %s %s from %s\n"

			ua := useragent.Parse(r.UserAgent())
			switch ua.Type {
			case 0:
				params = append(params, "unknown")
			case 1:
				params = append(params, "browser")
			case 2:
				params = append(params, "crawler")
			default:
				params = append(params, "other")
			}
			params = append(params, ua.Name)
			params = append(params, ua.Version)
			params = append(params, ua.OS)
			params = append(params, ua.OSVersion)
			device := "computer"
			if ua.Mobile {
				device = "mobile"
			}
			if ua.Tablet {
				device = "tablet"
			}
			params = append(params, device)
			text = text + "--- client: %s %s %s, OS: %s %s, device: %s\n"

			buf, _ := ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
			r.ParseForm()
			for i, v := range r.Form {
				text = text + "--- "
				text = text + i + " = "
				text = text + fmt.Sprint(v) + "\n"
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

			user := r.Header.Get("X_user")
			session := r.Header.Get("X_session")
			if user != "" || session != "" {
				text = text + "--- "
				if user != "" {
					text = text + "user: %s"
					params = append(params, user)
				}
				if session != "" {
					if user != "" {
						text = text + ", "
					}
					text = text + "session: %s"
					params = append(params, session)
				}
				text = text + "\n"
			}

			fmt.Printf(text, params...)

			lw := NewLoggingResponseWriter(w)
			a.Handler.ServeHTTP(lw, r)

			params = []interface{}{lw.statusCode}
			text = "<-- %d"
			user = w.Header().Get("X_user")
			session = w.Header().Get("X_session")
			if user != "" || session != "" {
				text = text + ", "
				if user != "" {
					text = text + "user: %s"
					params = append(params, user)
				}
				if session != "" {
					if user != "" {
						text = text + ", "
					}
					text = text + "session: %s"
					params = append(params, session)
				}
			}
			text = text + "\n\n"

			fmt.Printf(text, params...)
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
