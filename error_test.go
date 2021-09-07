package chassis_test

import (
	"fmt"
	"testing"

	"github.com/oligoden/chassis"

	"github.com/stretchr/testify/assert"
)

var err1 error = chassis.Mark("msg")
var err2 error = chassis.Mark("next", err1)

var err3 error = fmt.Errorf("msg")
var err4 error = chassis.Mark("next", err3)

func TestError(t *testing.T) {
	assert.Contains(t, err1.Error(), "msg")
	assert.Contains(t, err1.Error(), "error_test.go:12")
	assert.Contains(t, err2.Error(), "next")
	assert.Contains(t, err2.Error(), "error_test.go:13")
	assert.Contains(t, chassis.ErrorTrace(err2), "error_test.go:13")
	t.Error(chassis.ErrorTrace(err4))
}
