package view

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/oligoden/chassis/device/model"
)

type Operator interface {
	JSON(model.Operator)
	Error(model.Operator)
}

type Default struct {
	Response http.ResponseWriter `json:"-"`
	Err      error               `json:"-"`
}

func (v Default) JSON(m model.Operator) {
	if m.Err() != nil {
		log.Println(m.Err())
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	out, err := json.Marshal(m.Data())
	if err != nil {
		log.Println(m.Err())
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	v.Response.Write(out)
}

func (v Default) Error(m model.Operator) {
	if m.Err() != nil {
		log.Println(m.Err())
		if strings.Contains(m.Err().Error(), "bad request") {
			v.Response.WriteHeader(http.StatusBadRequest)
		} else {
			v.Response.WriteHeader(http.StatusInternalServerError)
		}
	}
}
