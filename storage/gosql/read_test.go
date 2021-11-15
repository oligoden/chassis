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
	exp += " LEFT JOIN record_users AS rus_testdata ON rus_testdata.record_id = testdata.hash"
	exp += " WHERE (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND rus_testdata.user_id = ?) OR testdata.owner_id = ?)"
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

func TestGenSelectGroup(t *testing.T) {
	data := TestData{}
	c := gosql.NewConnection(1, []uint{1, 2})
	c.GenSelect(data)
	q, vs := c.Query()

	exp := "SELECT testdata.* FROM testdata"
	exp += " LEFT JOIN record_groups AS rgs_testdata ON rgs_testdata.record_id = testdata.hash"
	exp += " LEFT JOIN record_users AS rus_testdata ON rus_testdata.record_id = testdata.hash"
	exp += " WHERE (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND rgs_testdata.group_id IN (?)) OR (testdata.perms LIKE ? AND rus_testdata.user_id = ?) OR testdata.owner_id = ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[%:%:%:%r% %:%:%r%:% %:%r%:%:% 1,2 %r%:%:%:% 1 1]"
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
	exp += " LEFT JOIN record_users AS rus_testdata ON rus_testdata.record_id = testdata.hash"
	exp += " WHERE test = ?"
	exp += " AND (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND rus_testdata.user_id = ?) OR testdata.owner_id = ?)"
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

func TestGenSelectTables(t *testing.T) {
	c := gosql.NewConnection(1, []uint{})
	join := gosql.NewJoin("LEFT JOIN sub_data on sub_data.test_data_id = testdata.id")
	c.AddModifiers(join)
	c.GenSelect(TestData{}, SubData{})
	q, vs := c.Query()

	exp := "SELECT testdata.*,subdata.* FROM testdata"
	exp += " LEFT JOIN sub_data on sub_data.test_data_id = testdata.id"
	exp += " LEFT JOIN record_users AS rus_testdata ON rus_testdata.record_id = testdata.hash"
	exp += " LEFT JOIN record_users AS rus_subdata ON rus_subdata.record_id = subdata.hash"
	exp += " WHERE (testdata.perms LIKE ? OR testdata.perms LIKE ? OR (testdata.perms LIKE ? AND rus_testdata.user_id = ?) OR testdata.owner_id = ?)"
	exp += " AND (subdata.perms LIKE ? OR subdata.perms LIKE ? OR (subdata.perms LIKE ? AND rus_subdata.user_id = ?) OR subdata.owner_id = ?)"
	got := q
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "[%:%:%:%r% %:%:%r%:% %r%:%:%:% 1 1 %:%:%:%r% %:%:%r%:% %r%:%:%:% 1 1]"
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

	q := "CREATE TABLE `testdata` ("
	q += "`field` varchar(255),"
	q += "`date` DATETIME NOT NULL DEFAULT '1000-01-01',"
	q += "`id` int unsigned AUTO_INCREMENT,"
	q += "`uc` varchar(255) UNIQUE,"
	q += "`owner_id` int unsigned,"
	q += "`perms` varchar(255),"
	q += "`hash` varchar(255),"
	q += "PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')")
	if err != nil {
		t.Fatal(err)
	}
	db.Exec("CREATE TABLE `subdata` (`test_data_id` int unsigned, `field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))")
	db.Exec("INSERT INTO `subdata` (`test_data_id`, `field`, `uc`, `owner_id`, `perms`, `hash`) VALUES (1, 'b', 'yy', 1, ':::', 'ijk')")

	s := gosql.New(uri)
	c := s.Connect(1, []uint{})
	join := gosql.NewJoin("JOIN `subdata` ON `subdata`.`test_data_id` = `testdata`.`id`")
	c.AddModifiers(join)
	e := TestData{}
	eSub := SubData{}
	c.Read(&e, &eSub)
	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "{a  1000-01-01 00:00:00 +0000 UTC [] [] {1 xx [] [] 1 ::: xyz}}"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "{1 b {1 yy [] [] 1 ::: ijk}}"
	got = fmt.Sprint(eSub)
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

	q := "CREATE TABLE `testdata` ("
	q += "`field` varchar(255),"
	q += "`date` DATETIME NOT NULL DEFAULT '1000-01-01',"
	q += "`id` int unsigned AUTO_INCREMENT,"
	q += "`uc` varchar(255) UNIQUE,"
	q += "`owner_id` int unsigned,"
	q += "`perms` varchar(255),"
	q += "`hash` varchar(255),"
	q += "PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}
	db.Exec("INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')")
	db.Exec("INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jkl')")
	db.Exec("CREATE TABLE `subdata` (`test_data_id` int unsigned, `field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))")
	db.Exec("INSERT INTO `subdata` (`test_data_id`, `field`, `uc`, `owner_id`, `perms`, `hash`) VALUES (1, 'c', 'zz', 1, ':::', 'cvb')")
	db.Exec("INSERT INTO `subdata` (`test_data_id`, `field`, `uc`, `owner_id`, `perms`, `hash`) VALUES (2, 'd', 'ff', 1, ':::', 'rty')")

	s := gosql.New(uri)
	c := s.Connect(1, []uint{})
	join := gosql.NewJoin("JOIN `subdata` ON `subdata`.`test_data_id` = `testdata`.`id`")
	c.AddModifiers(join)
	eSub := SubDataMap{}
	e := TestDataMap{}
	c.Read(e, eSub)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "map[xx:{a  1000-01-01 00:00:00 +0000 UTC [] [] {1 xx [] [] 1 ::: xyz}} yy:{b  1000-01-01 00:00:00 +0000 UTC [] [] {2 yy [] [] 1 ::: jkl}}]"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "map[ff:{2 d {2 ff [] [] 1 ::: rty}} zz:{1 c {1 zz [] [] 1 ::: cvb}}]"
	got = fmt.Sprint(eSub)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestReadMapID(t *testing.T) {
	testCleanup(t)

	db, err := sql.Open(dbt, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	q := "CREATE TABLE `testdata` ("
	q += "`field` varchar(255),"
	q += "`date` DATETIME NOT NULL DEFAULT '1000-01-01',"
	q += "`id` int unsigned AUTO_INCREMENT,"
	q += "`uc` varchar(255) UNIQUE,"
	q += "`owner_id` int unsigned,"
	q += "`perms` varchar(255),"
	q += "`hash` varchar(255),"
	q += "PRIMARY KEY (`id`))"
	_, err = db.Exec(q)
	if err != nil {
		t.Fatal(err)
	}
	db.Exec("INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('a', 'xx', 1, ':::', 'xyz')")
	db.Exec("INSERT INTO `testdata` (`field`, `uc`, `owner_id`, `perms`, `hash`) VALUES ('b', 'yy', 1, ':::', 'jkl')")
	db.Exec("CREATE TABLE `subdata` (`test_data_id` int unsigned, `field` varchar(255), `id` int unsigned AUTO_INCREMENT, `uc` varchar(255) UNIQUE, `owner_id` int unsigned, `perms` varchar(255), `hash` varchar(255), PRIMARY KEY (`id`))")
	db.Exec("INSERT INTO `subdata` (`test_data_id`, `field`, `uc`, `owner_id`, `perms`, `hash`) VALUES (1, 'c', 'zz', 1, ':::', 'cvb')")
	db.Exec("INSERT INTO `subdata` (`test_data_id`, `field`, `uc`, `owner_id`, `perms`, `hash`) VALUES (2, 'd', 'ff', 1, ':::', 'rty')")

	s := gosql.New(uri)
	c := s.Connect(1, []uint{})
	join := gosql.NewJoin("JOIN `subdata` ON `subdata`.`test_data_id` = `testdata`.`id`")
	c.AddModifiers(join)
	eSub := SubDataMapID{}
	e := TestDataMap{}
	c.Read(e, eSub)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "map[xx:{a  1000-01-01 00:00:00 +0000 UTC [] [] {1 xx [] [] 1 ::: xyz}} yy:{b  1000-01-01 00:00:00 +0000 UTC [] [] {2 yy [] [] 1 ::: jkl}}]"
	got := fmt.Sprint(e)
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "map[1:{1 c {1 zz [] [] 1 ::: cvb}} 2:{2 d {2 ff [] [] 1 ::: rty}}]"
	got = fmt.Sprint(eSub)
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

	q := "CREATE TABLE `testdata` ("
	q += "`field` varchar(255),"
	q += "`date` DATETIME NOT NULL DEFAULT '1000-01-01',"
	q += "`id` int unsigned AUTO_INCREMENT,"
	q += "`uc` varchar(255) UNIQUE,"
	q += "`owner_id` int unsigned,"
	q += "`perms` varchar(255),"
	q += "`hash` varchar(255),"
	q += "PRIMARY KEY (`id`))"
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

	s := gosql.New(uri)
	c := s.Connect(1, []uint{})
	e := TestDataSlice{}
	c.Read(&e)

	if c.Err() != nil {
		t.Error(c.Err())
	}

	exp := "[{a  1000-01-01 00:00:00 +0000 UTC [] [] {1 xx [] [] 1 ::: xyz}} {b  1000-01-01 00:00:00 +0000 UTC [] [] {2 yy [] [] 1 ::: jkl}}]"
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
	s := gosql.New(uri)

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

	exp := "usr    [] [] :::r vbn"
	got := fmt.Sprint(e)
	if !strings.Contains(got, exp) {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
