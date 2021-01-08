package gosql_test

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/oligoden/chassis/storage/gosql"
)

func TestGenSelect(t *testing.T) {
	data := TestData{}
	c := gosql.NewConnection(1, []uint{})
	c.GenSelect(data)
	q, vs := c.Query()

	exp := "SELECT testdata.* FROM testdata"
	exp += " LEFT JOIN record_groups on record_groups.record_id = testdata.hash"
	exp += " LEFT JOIN record_users on record_users.record_id = testdata.hash"
	exp += " WHERE (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND record_users.user_id = ?) OR testdata.owner_id = ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[%:%:%:%r% %:%:%r%:% %r%:%:%:% 1 1]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestGenSelectWhere(t *testing.T) {
	data := TestData{}
	c := gosql.NewConnection(1, []uint{})
	where := gosql.NewWhere("test = ?", 4)
	join1 := gosql.NewJoin("LEFT JOIN test_a on test_a.id = testdata.test_a_id")
	join2 := gosql.NewJoin("LEFT JOIN test_b on test_b.id = test_a.test_b_id")
	c.AddModifiers(join1, join2, where)
	c.GenSelect(data)
	q, vs := c.Query()

	exp := "SELECT testdata.* FROM testdata"
	exp += " LEFT JOIN test_a on test_a.id = testdata.test_a_id"
	exp += " LEFT JOIN test_b on test_b.id = test_a.test_b_id"
	exp += " LEFT JOIN record_groups on record_groups.record_id = testdata.hash"
	exp += " LEFT JOIN record_users on record_users.record_id = testdata.hash"
	exp += " WHERE test = ?"
	exp += " AND (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND record_users.user_id = ?) OR testdata.owner_id = ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[4 %:%:%:%r% %:%:%r%:% %r%:%:%:% 1 1]"
	got = fmt.Sprintf("%v", vs)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReadRecord(t *testing.T) {
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
	e := TestData{}
	c.Read(&e)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "{a  [] [] {1 xx [] [] 1 ::: xyz}}"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReadMap(t *testing.T) {
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

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jkl')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	s := gosql.New(dbt, uri)
	c := s.Connect(1, []uint{})
	e := TestDataMap{}
	c.Read(e)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "map[xx:{a  [] [] {1 xx [] [] 1 ::: xyz}} yy:{b  [] [] {2 yy [] [] 1 ::: jkl}}]"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReadSlice(t *testing.T) {
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

	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jkl')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	s := gosql.New(dbt, uri)
	c := s.Connect(1, []uint{})
	e := TestDataSlice{}
	c.Read(&e)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "[{a  [] [] {1 xx [] [] 1 ::: xyz}} {b  [] [] {2 yy [] [] 1 ::: jkl}}]"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	c.Read(e)
	if c.Err() == nil {
		t.Error("expected error")
	}
}

func TestReadUser(t *testing.T) {
	testCleanup(t)
	s := gosql.New(dbt, uri)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "INSERT INTO `users` (`uc`, `username`, `perms`, `hash`) VALUES ('c', 'usr', ':::r', 'vbn')"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}

	c := s.Connect(1, []uint{})
	e := gosql.UserRecords{}
	c.Read(&e)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "usr    [] [] :::r vbn}]"
	got := fmt.Sprint(e)
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
