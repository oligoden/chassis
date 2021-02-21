package device

import (
	"net/http"
)

func (d Default) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Bind()
		m.Create()

		v := d.NewView(w)
		v.JSON(m)
	})
}

func (d Default) Read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Read()

		v := d.NewView(w)
		v.JSON(m)
	})
}

func (d Default) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Read()
		m.Bind()
		m.Update()

		d.NewView(w).Error(m)
	})
}
