package model

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

func (m *Default) bind() error {
	rgx, _ := regexp.Compile("^/api/v[0-9]+/[a-z]+(/?|/(([a-zA-Z0-9]+)(/?|/.*)))$")
	matches := rgx.FindStringSubmatch(m.Request.URL.Path)

	if len(matches) == 0 {
		return fmt.Errorf("bad request, incorrect URL structure")
	}

	if matches[3] != "" {
		m.Data().UniqueCode(matches[2])
	}

	t := reflect.TypeOf(m.data).Elem()
	if t.Kind() == reflect.Struct {
		return m.structBind(m.data)
	}
	return nil
}

func (m *Default) structBind(s interface{}) error {
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()

	if t.Kind() != reflect.Struct {
		return fmt.Errorf("not a struct")
	}

	fmt.Println("\nBinding:")
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		if ft.Type.Name() == "RecordDefault" {
			err := m.structBind(fv.Addr().Interface())
			if err != nil {
				return err
			}
		}

		if tag, got := ft.Tag.Lookup("form"); got {
			if val := m.Request.FormValue(tag); val != "" {
				fmt.Printf("%d. %v (%v, %v), tag: '%v'\n", i+1, ft.Name, ft.Type.Name(), ft.Type.Kind(), ft.Tag.Get("form"))
				setType(ft.Type.Kind(), val, fv)
			}
		}
	}
	return nil
}

func setType(kind reflect.Kind, val string, fld reflect.Value) error {

	switch kind {
	case reflect.Ptr:
		return setType(fld.Elem().Kind(), val, fld.Elem())
	case reflect.Int:
		return setIntField(val, 0, fld)
	case reflect.Int8:
		return setIntField(val, 8, fld)
	case reflect.Int16:
		return setIntField(val, 16, fld)
	case reflect.Int32:
		return setIntField(val, 32, fld)
	case reflect.Int64:
		return setIntField(val, 64, fld)
	case reflect.Uint:
		return setUintField(val, 0, fld)
	case reflect.Uint8:
		return setUintField(val, 8, fld)
	case reflect.Uint16:
		return setUintField(val, 16, fld)
	case reflect.Uint32:
		return setUintField(val, 32, fld)
	case reflect.Uint64:
		return setUintField(val, 64, fld)
	case reflect.Bool:
		return setBoolField(val, fld)
	case reflect.Float32:
		return setFloatField(val, 32, fld)
	case reflect.Float64:
		return setFloatField(val, 64, fld)
	case reflect.String:
		fld.SetString(val)
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setIntField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	intVal, err := strconv.ParseInt(value, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0"
	}
	uintVal, err := strconv.ParseUint(value, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(value string, field reflect.Value) error {
	if value == "" {
		value = "false"
	}
	boolVal, err := strconv.ParseBool(value)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func setFloatField(value string, bitSize int, field reflect.Value) error {
	if value == "" {
		value = "0.0"
	}
	floatVal, err := strconv.ParseFloat(value, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}
