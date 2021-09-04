package model

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func (m *Default) bind() error {
	rgx, _ := regexp.Compile("^/[a-z]+(/?|/(([a-zA-Z0-9]+)(/?|/.*)))$")
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

	fmt.Println("\nbinding...")
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
			fmt.Printf("%s, tag: %s", ft.Name, tag)
			if val := m.Request.FormValue(tag); val != "" {
				fmt.Printf(", val: %s", val)
				setType(ft.Type, val, fv)
			}
			fmt.Println()
		}
	}
	return nil
}

func setType(ft reflect.Type, val string, fv reflect.Value) error {
	kind := ft.Kind()

	switch kind {
	case reflect.Ptr:
		return setType(fv.Elem().Type(), val, fv.Elem())
	case reflect.Int:
		return setIntField(val, 0, fv)
	case reflect.Int8:
		return setIntField(val, 8, fv)
	case reflect.Int16:
		return setIntField(val, 16, fv)
	case reflect.Int32:
		return setIntField(val, 32, fv)
	case reflect.Int64:
		return setIntField(val, 64, fv)
	case reflect.Uint:
		return setUintField(val, 0, fv)
	case reflect.Uint8:
		return setUintField(val, 8, fv)
	case reflect.Uint16:
		return setUintField(val, 16, fv)
	case reflect.Uint32:
		return setUintField(val, 32, fv)
	case reflect.Uint64:
		return setUintField(val, 64, fv)
	case reflect.Bool:
		return setBoolField(val, fv)
	case reflect.Float32:
		return setFloatField(val, 32, fv)
	case reflect.Float64:
		return setFloatField(val, 64, fv)
	case reflect.String:
		fv.SetString(val)
	case reflect.Struct:
		if ft.Name() == "Time" {
			return setTimeField(val, fv)
		}
		return errors.New("unknown struct type")
	default:
		return errors.New("unknown type")
	}
	return nil
}

func setTimeField(value string, field reflect.Value) error {
	if value == "" {
		value = "1000-01-01 00:00:00"
	}

	loc, err := time.LoadLocation("Local")
	if err != nil {
		return err
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05.999999", value, loc)
	if err != nil {
		t, err = time.ParseInLocation("2006-01-02", value, loc)
		if err != nil {
			return err
		}
	}

	field.Set(reflect.ValueOf(t))
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
