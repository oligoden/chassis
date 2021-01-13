package gosql_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestGenDelete(t *testing.T) {
	e := &TestData{
		Field: "test",
	}
	e.ID = 1

	c := gosql.NewConnection(1, []uint{})
	where := gosql.NewWhere("id = ?", 1)
	c.AddModifiers(where)
	c.GenDelete(e)
	q, vs := c.Query()

	exp := "DELETE testdata.* FROM testdata LEFT JOIN record_groups on record_groups.record_id = testdata.hash LEFT JOIN record_users on record_users.record_id = testdata.hash WHERE id = ? AND (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND record_users.user_id = ?) OR testdata.owner_id = ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[1 %:%:%:%d% %:%:%d%:% %d%:%:%:% 1 1]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestDelete(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` (`field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	s := gosql.New(dbt, uri)
	c := s.Connect(1, []uint{})

	e := &TestData{}
	c.Delete(e)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	var field string
	err = db.QueryRow("SELECT field from testdata").Scan(&field)
	if err == nil {
		t.Error("expected error")
	}
}
