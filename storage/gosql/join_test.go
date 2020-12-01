package gosql_test

import (
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestJoins(t *testing.T) {
	j := gosql.NewJoin("LEFT JOIN b on b.ba = a.aa")
	j.Add("LEFT JOIN c on c.ca = a.aa")
	q := j.Compile()

	exp := " LEFT JOIN b on b.ba = a.aa LEFT JOIN c on c.ca = a.aa"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
