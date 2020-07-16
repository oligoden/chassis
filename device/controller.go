package device

import (
	"net/http"
)

func (d Default) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Bind()
		v := d.NewView(w)
		db := d.Store.CreateDB(m.User())

		m.Create(db)
		db.Close()

		v.JSON(m)
	})
}

func (d Default) Read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		m.Bind()
		v := d.NewView(w)
		db := d.Store.ReadDB(m.User())

		m.Read(db)
		db.Close()

		v.JSON(m)
	})
}

func (d Default) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := d.NewModel(r)
		v := d.NewView(w)
		dbRead := d.Store.ReadDB(m.User())

		m.Read(dbRead, "with-update")
		m.Bind()
		db := dbRead.ReaderToUpdater()
		m.Update(db)
		dbRead.Close()

		v.JSON(m)
	})
}
