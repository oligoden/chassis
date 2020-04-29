package data_test

import (
	"encoding/json"
	"testing"

	"github.com/oligoden/chassis/device/model/data"
)

func TestPrepare(t *testing.T) {
	x := data.Default{}
	if x.Prepare() != nil {
		t.Error("expected nil")
	}
}

func TestRead(t *testing.T) {
	x := data.Default{}
	if x.Read(nil) != nil {
		t.Error("expected nil")
	}
}

func TestComplete(t *testing.T) {
	x := data.Default{}
	if x.Complete() != nil {
		t.Error("expected nil")
	}
}

func TestResponse(t *testing.T) {
	x := data.Default{
		UC: "abc",
	}
	gotJSON, _ := json.Marshal(x.Response())
	got := string(gotJSON)
	exp := `{"uc":"abc"}`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestTableName(t *testing.T) {
	x := data.Default{}
	got := x.TableName()
	exp := `models`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestUniqueCode(t *testing.T) {
	x := data.Default{}

	got := x.UniqueCode()
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = x.UniqueCode("abc")
	exp = `abc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = x.UniqueCode()
	exp = `abc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestPermissions(t *testing.T) {
	x := data.Default{}

	got := x.Permissions()
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = x.Permissions(":::")
	exp = `:::`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	got = x.Permissions()
	exp = `:::`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}

func TestOwner(t *testing.T) {
	x := data.Default{}

	got := x.Owner()
	exp := uint(0)
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = x.Owner(1)
	exp = 1
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = x.Owner()
	exp = 1
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}
}

func TestGroups(t *testing.T) {
	x := data.Default{}

	got := len(x.Groups())
	exp := len([]uint{})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = len(x.Groups(1))
	exp = len([]uint{1})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}

	got = len(x.Groups(2, 3))
	exp = len([]uint{1, 2, 3})
	if got != exp {
		t.Errorf(`expected '%d', got '%d'`, exp, got)
	}
}

func TestHash(t *testing.T) {
	x := data.Default{}

	got := x.Hash
	exp := ``
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}

	err := x.Hasher()
	if err != nil {
		t.Error(err)
	}

	got = x.Hash
	exp = `a50f26ed4478d16faec97103708fa99883adafcc`
	if got != exp {
		t.Errorf(`expected '%s', got '%s'`, exp, got)
	}
}
