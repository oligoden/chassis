package adapter

import (
	"net/http"

	"github.com/oligoden/chassis/storage/gosql"
)

type Mux struct {
	Mux    *http.ServeMux
	Domain string
	Stores map[string]*gosql.Store
	RPDs   []string
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
