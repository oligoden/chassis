package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestReadWhere(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &TestModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	db.Close()

	storage := gormdb.New(dbt, uri)
	dbRead := storage.ReadDB(0, []uint{})
	m = &TestModel{}
	dbRead.Where("testmodels.test_model_id = ?", 1).First(m)
	dbRead.Close()
	if dbRead.Error() != nil {
		t.Error(dbRead.Error())
	}

	exp := uint(1)
	got := m.TestModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}

func TestReadNewRecord(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &TestModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	db.Close()

	storage := gormdb.New(dbt, uri)
	dbRead := storage.ReadDB(0, []uint{})
	m = &TestModel{}
	if !dbRead.NewRecord(m) {
		t.Errorf(`expected new record`)
	}
	dbRead.First(m)
	if dbRead.NewRecord(m) {
		t.Errorf(`expected existing record`)
	}
	dbRead.Close()
	if dbRead.Error() != nil {
		t.Error(dbRead.Error())
	}

	exp := uint(1)
	got := m.TestModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}

func TestReadFirstWithError(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &TestModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	db.Close()

	// simulate error
	storage := gormdb.New("", "")
	dbRead := storage.ReadDB(0, []uint{})

	if dbRead.Error() == nil {
		t.Error(`expected error`)
	}

	m = &TestModel{}
	dbRead.First(m)
	dbRead.Close()

	exp := uint(0)
	got := m.TestModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}

func TestReadFirst(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})
	setupDBTable(&SubModel{})
	setupDBTable(&WeakModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &SubModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	mWeakModel := &WeakModel{Field: "a", Perms: ":::r"}
	db.Create(mWeakModel)
	db.Close()

	store := gormdb.New(dbt, uri)
	dbRead := store.ReadDB(0, []uint{})
	m = &SubModel{}
	dbRead.First(m)
	if dbRead.Error() != nil {
		t.Error(dbRead.Error())
	}
	mWeakModel = &WeakModel{}
	dbRead.First(mWeakModel, "weak_models")
	if dbRead.Error() != nil {
		t.Error(dbRead.Error())
	}

	exp := uint(1)
	got := m.SubModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
	exp = uint(1)
	got = mWeakModel.WeakModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}

	dbRead.First(mWeakModel)
	expErr := "model is not assertable as an table namer"
	gotErr := dbRead.Error().Error()
	if expErr != gotErr {
		t.Errorf(`expected "%s", got "%s"`, expErr, gotErr)
	}
	dbRead.Close()
}

func TestReadPreloadWithError(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})
	setupDBTable(&SubModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &TestModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	mSubModel := &SubModel{
		UC:          "asd",
		TestModelID: 1,
		Field:       "a",
		Perms:       ":::r",
	}
	db.Create(mSubModel)
	mSubModel = &SubModel{
		UC:          "gfd",
		TestModelID: 1,
		Field:       "a",
		Perms:       ":::r",
	}
	db.Create(mSubModel)
	db.Close()

	// simulate error
	storage := gormdb.New("", "")
	dbRead := storage.ReadDB(0, []uint{})

	if dbRead.Error() == nil {
		t.Error(`expected error`)
	}

	m = &TestModel{}
	dbRead.Preload("SubModels", "submodels").First(m)
	dbRead.Close()

	if len(m.SubModels) != 0 {
		t.Error(`expected no preloaded submodels`)
	}
}

func TestReadPreload(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})
	setupDBTable(&SubModel{})

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)
	m := &TestModel{Field: "a", Perms: ":::r"}
	db.Create(m)
	mSubModel := &SubModel{
		UC:          "asd",
		TestModelID: 1,
		Field:       "a",
		Perms:       ":::r",
	}
	db.Create(mSubModel)
	mSubModel = &SubModel{
		UC:          "gfd",
		TestModelID: 1,
		Field:       "a",
		Perms:       ":::r",
	}
	db.Create(mSubModel)
	db.Close()

	storage := gormdb.New(dbt, uri)
	dbRead := storage.ReadDB(0, []uint{})
	m = &TestModel{}
	dbRead.Preload("SubModels", "submodels").First(m)
	dbRead.Close()
	if dbRead.Error() != nil {
		t.Error(dbRead.Error())
	}

	if len(m.SubModels) == 0 {
		t.Error(`expected preloaded submodels`)
	}
	exp := uint(1)
	got := m.TestModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
	exp = uint(1)
	got = m.SubModels[0].SubModelID
	if exp != got {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}
