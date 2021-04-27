package chassis_test

import (
	"regexp"
	"testing"

	"github.com/oligoden/chassis"
)

var err error = chassis.Mark("msg")
var err2 error = chassis.Mark("next", err)

func TestSingleError(t *testing.T) {
	exp := "(?s)msg (.*error_test.go:10)"
	got := err.Error()
	if m, err := regexp.MatchString(exp, got); !m || err != nil {
		if err != nil {
			t.Error(err)
		}
		if !m {
			t.Errorf("mismatch\nexpected %s\ngot %s", exp, got)
		}
	}
}

func TestStackError(t *testing.T) {
	exp := "(?s)next (.*error_test.go:11)"
	got := err2.Error()
	if m, err := regexp.MatchString(exp, got); !m || err != nil {
		if err != nil {
			t.Error(err)
		}
		if !m {
			t.Errorf("mismatch\nexpected %s\ngot %s", exp, got)
		}
	}

	exp = "(?s)msg (.*error_test.go:10)..next (.*error_test.go:11)"
	got = chassis.ErrorTrace(err2)
	if m, err := regexp.MatchString(exp, got); !m || err != nil {
		if err != nil {
			t.Error(err)
		}
		if !m {
			t.Errorf("mismatch\nexpected %s\ngot\n%s", exp, got)
		}
	}
}
