package gosql_test

import (
	"database/sql"
	"fmt"
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

func TestReadRecord(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
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

	exp := "{1 a  [] [] {xx [] [] 1 ::: xyz}}"
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

	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
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

	exp := "map[xx:{1 a  [] [] {xx [] [] 1 ::: xyz}} yy:{2 b  [] [] {yy [] [] 1 ::: jkl}}]"
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

	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
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

	exp := "[{1 a  [] [] {xx [] [] 1 ::: xyz}} {2 b  [] [] {yy [] [] 1 ::: jkl}}]"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

// func TestReadMapWhere(t *testing.T) {
// 	testCleanup(t)

// 	db, err := sql.Open(dbt, uri)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer db.Close()

// 	q := "CREATE TABLE `testdata` (`id` int unsigned AUTO_INCREMENT, `field` varchar(255), `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))"
// 	_, err = db.Exec(q)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')"
// 	_, err = db.Exec(q)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	q = "INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jfk')"
// 	_, err = db.Exec(q)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	s := gosql.New(dbt, uri)
// 	c := s.Connect(1, []uint{})
// 	e := TestDataMap{}
// 	c.Where("id = ?", 0)
// 	c.Read(e)

// 	if s.Err() != nil {
// 		t.Error(s.Err())
// 	}

// 	exp := "map[xx:{1 a  [] [] xx [] [] 1 ::: xyz}]"
// 	got := fmt.Sprint(e)
// 	if exp != got {
// 		t.Errorf(`expected "%s", got "%s"`, exp, got)
// 	}
// }
