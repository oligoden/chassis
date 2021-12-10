package adapter

import (
	"net/http"
	"net/url"

	"github.com/oligoden/chassis/storage/gosql"
)

type Mux struct {
	Mux    *http.ServeMux
	Domain string
	URL    *url.URL
	Stores map[string]*gosql.Store
	RPDs   []string
	Err    error
}

func NewMux() *Mux {
	return &Mux{
		Mux:    http.NewServeMux(),
		Stores: map[string]*gosql.Store{},
		RPDs:   []string{},
	}
}

func (mx *Mux) SetDomain(domain string) *Mux {
	mx.Domain = domain
	return mx
}

func (mx *Mux) SetURL(us string) *Mux {
	u, err := url.Parse(us)
	if err != nil {
		mx.Err = err
		return mx
	}

	mx.URL = u
	return mx
}

func (mx *Mux) SetStore(key string, store *gosql.Store) *Mux {
	mx.Stores[key] = store
	return mx
}

func (mx *Mux) AddRPD(dest string) *Mux {
	mx.RPDs = append(mx.RPDs, dest)
	return mx
}

func (mx *Mux) Compile(hs func(*Mux)) *http.ServeMux {
	hs(mx)
	return mx.Mux
}

func (mx *Mux) ServeMux() *http.ServeMux {
	return mx.Mux
}

func (mx *Mux) Adapter() Adapter {
	return Adapter{
		Host: mx.URL.Hostname(),
		mx:   mx,
	}
}

func (mx *Mux) Handle(pattern string) Adapter {
	return Adapter{
		Host:    mx.URL.Hostname(),
		pattern: pattern,
		mx:      mx,
	}
}
