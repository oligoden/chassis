package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestUpdateWithError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::u"}
	db.Create(m)

	// simulate error with incorrect connection
	storage := gormdb.New("", "")
	dbUpdate := storage.UpdateDB(2, []uint{})

	if dbUpdate.Error() == nil {
		t.Error(`expected error`)
	}

	m.Field = "b"
	dbUpdate.Save(m)
	dbUpdate.Close()

	m = &TestModel{}
	db.First(m)
	db.Close()

	exp := "a"
	got := m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUpdateAuthFailure(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbUpdate := storage.UpdateDB(2, []uint{})
	m.Field = "b"
	dbUpdate.Save(m)
	dbUpdate.Close()

	if dbUpdate.Error() == nil {
		t.Error(`expected error`)
		t.FailNow()
	}
	exp := "update authorization failed"
	got := dbUpdate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m = &TestModel{}
	db.First(m)
	db.Close()

	exp = "a"
	got = m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUpdateAuthError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: "::"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbUpdate := storage.UpdateDB(2, []uint{})
	m.Field = "b"
	dbUpdate.Save(m)
	dbUpdate.Close()

	if dbUpdate.Error() == nil {
		t.Error(`expected error`)
		t.FailNow()
	}
	exp := "the model has incorrect permissions format"
	got := dbUpdate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m = &TestModel{}
	db.First(m)
	db.Close()

	exp = "a"
	got = m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUpdate(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::u"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbUpdate := storage.UpdateDB(2, []uint{})
	m.Field = "b"
	dbUpdate.Save(m)
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

func TestUpdateWithCreateFail(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::c"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbUpdate := storage.UpdateDB(2, []uint{})
	m.Field = "b"
	dbUpdate.Save(m)
	dbUpdate.Close()

	if dbUpdate.Error() == nil {
		t.Error(`expected error`)
		t.FailNow()
	}
	exp := "update authorization failed"
	got := dbUpdate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m = &TestModel{}
	db.First(m)
	db.Close()

	exp = "a"
	got = m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestUpdateWithCreate(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::c"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbUpdate := storage.UpdateDB(2, []uint{})
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
