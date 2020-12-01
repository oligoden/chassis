package gosql_test

import (
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestWhere(t *testing.T) {
	w := gosql.NewWhere("a = ?", 1)
	q, vs := w.Compile()

	exp := " WHERE a = ?"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[1]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestWhereAndOr(t *testing.T) {
	w := gosql.NewWhere("a = ?", 1)
	w.And("b > ?", 2)
	w.Or("c LIKE ?", "abc")
	q, vs := w.Compile()

	exp := " WHERE a = ? AND b > ? OR c LIKE ?"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[1 2 abc]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestWhereAndOrGroup(t *testing.T) {
	w := gosql.NewWhere("a = ?", 1)
	wg := gosql.NewWhere("b > ?", 2)
	w.AndGroup(wg)
	wg = gosql.NewWhere("c LIKE ?", "abc")
	w.OrGroup(wg)
	q, vs := w.Compile()

	exp := " WHERE a = ? AND (b > ?) OR (c LIKE ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[1 2 abc]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
