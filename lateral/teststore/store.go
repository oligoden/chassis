package teststore

import (
	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/storage"
)

type CRUD struct{}

func Create(e data.Operator) {

}

func Read(e ...data.Operator) {

}

func Update(e data.Operator) {

}

func Delete(e data.Operator) {

}

func AddModifiers(mods ...storage.Modifier) {}

func Err(errs ...error) error { return nil }
