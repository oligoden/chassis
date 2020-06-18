package device

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/view"
	"github.com/oligoden/chassis/storage"
)

type Default struct {
	NewModel func(r *http.Request) model.Operator
	NewView  func(w http.ResponseWriter) view.Operator
	Store    storage.Storer
	user     uint
	groups   []uint
	err      error
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

func (d *Default) Bind(r *http.Request) {
	if d.err != nil {
		return
	}

	u := r.Header.Get("X_Session_User")
	user, _ := strconv.Atoi(u)
	d.user = uint(user)

	for _, g := range strings.Split(r.Header.Get("X_User_Groups"), ",") {
		group, _ := strconv.Atoi(g)
		d.groups = append(d.groups, uint(group))
	}
}

func (d Default) User() (uint, []uint) {
	return d.user, d.groups
}
