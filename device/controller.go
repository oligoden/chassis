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

// func (d Default) Update() func(w http.ResponseWriter, r *http.Request) error {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		d.Bind(r)
// 		m := d.NewModel(w, r)

// 		dbRead := d.Store.ReadDB(d.User())
// 		m.Read(dbRead)
// 		m.Bind()
// 		db := dbRead.ReaderToUpdater()
// 		m.Update(db)
// 		dbRead.Close()
// 		return m.Response()
// 	}
// }
