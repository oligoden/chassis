package device

import (
	"fmt"
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

		m.Bind()
		if m.Data().UniqueCode() == "" {
			m.NewData("list")
			fmt.Println("list")
		}
		m.Read()
		fmt.Printf("data %+v", m.Data())

		v := d.NewView(w)
		v.JSON(m)
	})
}

func (d Default) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)

		m.Bind()
		if m.Data().UniqueCode() != "" {
			m.Read()
		}

		m.Bind()
		m.Update()

		d.NewView(w).Error(m)
	})
}

func (d Default) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Bind()
		m.Delete()

		d.NewView(w).Error(m)
	})
}
