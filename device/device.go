package device

import (
	"log"
	"net/http"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage"
)

type Default struct {
	NewModel func(r *http.Request) model.Operator
	NewView  func(w http.ResponseWriter) view.Operator
	Store    storage.Storer
}

func NewDevice(nm func(r *http.Request) model.Operator, nv func(w http.ResponseWriter) view.Operator, store storage.Storer) Default {
	d := Default{}
	d.Store = store
	d.NewModel = nm
	d.NewView = nv
	return d
}

func (d Default) Manage(action string) {
	m := d.NewModel(nil)

	db := d.Store.ManageDB()
	if db.Error() != nil {
		log.Fatal(db.Error())
	}

	m.Manage(db, action)
	if m.Error() != nil {
		log.Fatal(m.Error())
	}
	db.Close()
}
