package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage/gormdb"
)

func TestUniqueCodeGeneration(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	db.Close()

	storage := gormdb.New(dbt, uri)
	storage.UniqueCodeFunc(chassis.RandNumberString)
	storage.UniqueCodeLength(1)

	dbCreate := storage.CreateDB(0, []uint{})
	m := &TestModel{}
	for i := 0; i < 15; i++ {
		m = &TestModel{}
		m.Perms = ":::c"
		dbCreate.Create(m)
	}
	dbCreate.Close()
	if dbCreate.Error() != nil {
		t.Error(dbCreate.Error())
	}

	if len(m.UC) <= 1 {
		t.Errorf(`expected "> 1", got "%d"`, len(m.UC))
	}
}

func TestCreateWithError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	db.Close()

	// simulate error with incorrect connection
	storage := gormdb.New("", "")
	dbCreate := storage.CreateDB(0, []uint{})

	if dbCreate.Error() == nil {
		t.Error(`expected error`)
	}

	m := &TestModel{}
	dbCreate.Create(m)
	dbCreate.Close()

	exp := uint(0)
	got := m.TestModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}

func TestCreateAuthFailure(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	db.Close()

	storage := gormdb.New(dbt, uri)
	dbCreate := storage.CreateDB(2, []uint{})

	m := &TestModel{}
	m.Perms = ":::"
	dbCreate.Create(m)
	dbCreate.Close()

	if dbCreate.Error() == nil {
		t.Error(`expected error`)
		t.FailNow()
	}
	exp := "create authorization failed"
	got := dbCreate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreateAuthError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	db.Close()

	storage := gormdb.New(dbt, uri)
	dbCreate := storage.CreateDB(0, []uint{})

	m := &TestModel{}
	m.Perms = "::"
	dbCreate.Create(m)
	dbCreate.Close()

	if dbCreate.Error() == nil {
		t.Error(`expected error`)
		t.FailNow()
	}
	exp := "the model has incorrect permissions format"
	got := dbCreate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreateToUpdate(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	storage := gormdb.New(dbt, uri)
	dbCreate := storage.CreateDB(0, []uint{})
	if dbCreate.Error() != nil {
		t.Error(dbCreate.Error())
	}

	m := &TestModel{Field: "a", Perms: ":::c"}
	dbCreate.Create(m)

	dbUpdate := dbCreate.CreaterToUpdater()
	m.Field = "b"
	dbUpdate.Save(m, "with-create")
	dbUpdate.Close()
	if dbUpdate.Error() != nil {
		t.Error(dbUpdate.Error())
	}

	m = &TestModel{}
	db.First(m)
	db.Close()

	exp := "b"
	got := m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}
