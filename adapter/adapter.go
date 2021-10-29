package adapter

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	useragent "github.com/mssola/user_agent"
)

type Adapter struct {
	Handler http.Handler
	Host    string
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

func New(host string) Adapter {
	return Adapter{
		Host: host,
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

func MethodNotAllowed() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}),
	}
}

func NotFound() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
	}
}

func (a Adapter) Core(h http.Handler) Adapter {
	return Adapter{
		Handler: h,
		Host:    a.Host,
	}
}

func (a Adapter) MNA() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}),
		Host: a.Host,
	}
}

func (a Adapter) MethodNotAllowed() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}),
		Host: a.Host,
	}
}

func (a Adapter) NotFound() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}),
		Host: a.Host,
	}
}

func (a Adapter) Notify() Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			params := []interface{}{r.Method, r.URL.String(), r.RemoteAddr}
			text := "--> %s %s from %s\n"

			ua := useragent.New(r.UserAgent())

			if ua != nil {
				if ua.Bot() {
					params = append(params, "bot")
				} else {
					params = append(params, "browser")
				}

				name, version := ua.Browser()
				params = append(params, name)
				params = append(params, version)
				params = append(params, ua.Platform())
				params = append(params, ua.OS())

				if ua.Mobile() {
					params = append(params, "mobile")
				} else {
					params = append(params, "desktop")
				}

				text = text + "--- client: %s %s %s, OS: %s %s, device: %s\n"
			} else {
				text = text + "--- client: unknown\n"
			}

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
		Host: a.Host,
	}
}

type doneWriter struct {
	http.ResponseWriter
	done bool
}

func (w *doneWriter) WriteHeader(status int) {
	w.done = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *doneWriter) Write(b []byte) (int, error) {
	w.done = true
	return w.ResponseWriter.Write(b)
}

func (a Adapter) And(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dw := &doneWriter{ResponseWriter: w}
			h.ServeHTTP(dw, r)
			if dw.done {
				return
			}
			a.Handler.ServeHTTP(w, r)
		}),
		Host: a.Host,
	}
}

func (a Adapter) Get(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				fmt.Println("--- executing GET handler")
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
		Host: a.Host,
	}
}

func (a Adapter) Post(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				fmt.Println("--- executing POST handler")
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
		Host: a.Host,
	}
}

func (a Adapter) Put(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				fmt.Println("--- executing PUT handler")
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
		Host: a.Host,
	}
}

func (a Adapter) Delete(h http.Handler) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodDelete {
				fmt.Println("--- executing DELETE handler")
				h.ServeHTTP(w, r)
			} else {
				a.Handler.ServeHTTP(w, r)
			}
		}),
		Host: a.Host,
	}
}

func (a Adapter) SubDomain(h http.Handler, rules ...string) Adapter {
	return Adapter{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host == a.Host {
				a.Handler.ServeHTTP(w, r)
				return
			}

			fmt.Println("requested host", r.Host, "on domain", a.Host)
			subdomain := strings.TrimSuffix(r.Host, a.Host)
			subdomain = strings.TrimSuffix(subdomain, ".")

			for _, rule := range rules {
				if rule == "" {
					continue
				}

				if strings.HasPrefix(rule, "-") && subdomain == rule[1:] {
					fmt.Println("ignoring by rule", rule)
					a.Handler.ServeHTTP(w, r)
					return
				}
			}

			for _, rule := range rules {
				if rule == "" {
					continue
				}

				if !strings.HasPrefix(rule, "-") && subdomain == rule {
					fmt.Println("rerouting by rule", rule)
					h.ServeHTTP(w, r)
					return
				}
			}

			h.ServeHTTP(w, r)
		}),
		Host: a.Host,
	}
}

func (a Adapter) Entry() http.Handler {
	return a.Handler
}
