package view

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/oligoden/chassis/device/model"
)

type Operator interface {
	JSON(model.Operator)
}

type Default struct {
	Response http.ResponseWriter `json:"-"`
	Err      error               `json:"-"`
}

func (v Default) JSON(m model.Operator) {
	if m.Error() != nil {
		log.Println(m.Error())
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(m.Data())
	if err != nil {
		log.Println(m.Error())
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	v.Response.Write(out)
}
