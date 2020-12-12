package data_test

import (
	"testing"

	"github.com/oligoden/chassis/device/model/data"
)

func TestPrepare(t *testing.T) {
	e := data.Default{}
	if e.Prepare() != nil {
		t.Error("expected nil")
	}
}

func TestComplete(t *testing.T) {
	e := data.Default{}
	if e.Complete() != nil {
		t.Error("expected nil")
	}
}

func TestTableName(t *testing.T) {
	e := data.Default{}
	got := e.TableName()
	exp := `models`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestUniqueCode(t *testing.T) {
	e := data.Default{}

	got := e.UniqueCode()
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = e.UniqueCode("abc")
	exp = `abc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = e.UniqueCode()
	exp = `abc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestPermissions(t *testing.T) {
	e := data.Default{}

	got := e.Permissions()
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = e.Permissions(":::")
	exp = `:::`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = e.Permissions()
	exp = `:::`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestOwner(t *testing.T) {
	e := data.Default{}

	got := e.Owner()
	exp := uint(0)
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = e.Owner(1)
	exp = 1
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = e.Owner()
	exp = 1
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}
}

func TestGroups(t *testing.T) {
	e := data.Default{}

	got := len(e.Groups())
	exp := len([]uint{})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = len(e.Groups(1))
	exp = len([]uint{1})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = len(e.Groups(2, 3))
	exp = len([]uint{1, 2, 3})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}
}

func TestHash(t *testing.T) {
	e := data.Default{}

	got := e.Hash
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	err := e.Hasher()
	if err != nil {
		t.Error(err)
	}

	got = e.Hash
	exp = `a50f26ed4478d16faec97103708fa99883adafcc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}
