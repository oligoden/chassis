package gormdb_test

import (
	"log"
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestMigrate(t *testing.T) {
	s := gormdb.New(dbt, uri)

	dbManage := s.ManageDB()
	dbManage.Manage(&TestModel{}, "dropIfExists")
	dbManage.Manage(&TestModel{}, "migrate")
	dbManage.Close()
	if dbManage.Error() != nil {
		t.Error(dbManage.Error())
	}

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	if !db.HasTable(&TestModel{}) {
		t.Error("table expected")
	}
	db.Close()
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}

func TestDrop(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	s := gormdb.New(dbt, uri)
	dbManage := s.ManageDB()
	dbManage.Manage(&TestModel{}, "drop")
	dbManage.Close()
	if dbManage.Error() != nil {
		t.Error(dbManage.Error())
	}

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	if db.HasTable(&TestModel{}) {
		t.Error("table not expected")
	}
	db.Close()
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}

func TestDropIfExist(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	s := gormdb.New(dbt, uri)
	dbManage := s.ManageDB()
	dbManage.Manage(&TestModel{}, "dropIfExists")
	dbManage.Close()
	if dbManage.Error() != nil {
		t.Error(dbManage.Error())
	}

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	if db.HasTable(&TestModel{}) {
		t.Error("table not expected")
	}
	db.Close()
	err = db.Error
	if err != nil {
		log.Fatal(err)
	}
}

func TestManageErrors(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	s := gormdb.New(dbt, "")
	dbManage := s.ManageDB()
	dbManage.Manage(&TestModel{}, "drop")
	dbManage.Close()
	if dbManage.Error() == nil {
		t.Error("expected error")
	}
}
