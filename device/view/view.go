package view

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/oligoden/chassis"
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
		v.Error(m)
		return
	}

	out, err := json.Marshal(m.Data())
	if err != nil {
		log.Println(err)
		http.Error(v.Response, "an error occured", http.StatusInternalServerError)
		return
	}

	v.Response.Write(out)
}

func (v Default) Error(m model.Operator) {
	if m.Err() != nil {
		if strings.Contains(m.Err().Error(), "bad request") {
			fmt.Println(m.Err())
			v.Response.WriteHeader(http.StatusBadRequest)
		} else {
			fmt.Printf("ERROR\n%s\n", chassis.ErrorTrace(m.Err()))
			v.Response.WriteHeader(http.StatusInternalServerError)
		}
	}
}
