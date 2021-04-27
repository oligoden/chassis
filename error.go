package chassis

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

type ChassisError struct {
	TS       time.Time
	Message  string
	Function string
	File     string
	Line     uint
	Child    error
}

func (e ChassisError) Error() string {
	return fmt.Sprintf("%s (%s:%d)", e.Message, e.File, e.Line)
}

func ErrorTrace(err error) string {
	e := err.Error()
	err = errors.Unwrap(err)
	for err != nil {
		e = fmt.Sprintf("%s\n%s", err.Error(), e)
		err = errors.Unwrap(err)
	}

	return e
}

func (e ChassisError) Unwrap() error {
	return e.Child
}

func Mark(msg string, errs ...error) error {
	e := ChassisError{}

	if len(errs) > 0 {
		e.Child = errs[0]
	}

	e.Message = msg
	pc, file, line, _ := runtime.Caller(1)
	e.Function = runtime.FuncForPC(pc).Name()
	e.File = file
	e.Line = uint(line)

	return e
}
