package device

import (
	"net/http"

	"github.com/oligoden/chassis/device/model"
	"github.com/oligoden/chassis/device/view"
)

type NewModelFunc func(r *http.Request) model.Operator
type NewViewFunc func(w http.ResponseWriter) view.Operator

type Default struct {
	NewModel NewModelFunc
	NewView  NewViewFunc
	Store    model.Connector
}

func NewDevice(nm NewModelFunc, nv NewViewFunc, store model.Connector) Default {
	d := Default{}
	d.Store = store
	d.NewModel = nm
	d.NewView = nv
	return d
}
