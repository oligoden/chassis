package gormdb_test

import (
	"fmt"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage/gormdb"
)

func TestAssociateAppendWithError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	setupDBTable(&SubModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::"}
	db.Create(m)

	// simulate error with incorrect connection
	storage := gormdb.New("", "")
	dbAssociate := storage.AssociateDB(0, []uint{})

	if dbAssociate.Error() == nil {
		t.Error(`expected error`)
	}

	mSub := &SubModel{}
	dbAssociate.Append("SubModels", m, mSub)
	dbAssociate.Close()

	m = &TestModel{}
	db.Preload("SubModels").First(m)
	db.Close()

	exp := 0
	got := len(m.SubModels)
	if got != exp {
		t.Errorf(`expected "%d", got "%d"`, exp, got)
	}
}

func TestAssociateAppendCreateAuthFailure(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	setupDBTable(&SubModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbAssociate := storage.AssociateDB(0, []uint{})
	mSub := &SubModel{Perms: ":::"}
	dbAssociate.Append("SubModels", m, mSub)
	dbAssociate.Close()

	if dbAssociate.Error() == nil {
		t.Error(`expected error`)
	}
	exp := "associate create authorization failed"
	got := dbAssociate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m = &TestModel{}
	db.Preload("SubModels").First(m)
	db.Close()

	expInt := 0
	gotInt := len(m.SubModels)
	if gotInt != expInt {
		t.Errorf(`expected "%d", got "%d"`, expInt, gotInt)
	}
}

func TestAssociateAppendCreateAuthError(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	setupDBTable(&SubModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbAssociate := storage.AssociateDB(0, []uint{})
	mSub := &SubModel{Perms: "::"}
	dbAssociate.Append("SubModels", m, mSub)
	dbAssociate.Close()

	if dbAssociate.Error() == nil {
		t.Error(`expected error`)
	}
	exp := "the model has incorrect permissions format"
	got := dbAssociate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	m = &TestModel{}
	db.Preload("SubModels").First(m)
	db.Close()

	expInt := 0
	gotInt := len(m.SubModels)
	if gotInt != expInt {
		t.Errorf(`expected "%d", got "%d"`, expInt, gotInt)
	}
}

func TestAssociateAppendCreate(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	setupDBTable(&SubModel{}, db)

	m := &TestModel{Field: "a", Perms: ":::"}
	db.Create(m)

	storage := gormdb.New(dbt, uri)
	dbAssociate := storage.AssociateDB(1, []uint{})
	mSub := &SubModel{Perms: ":::c"}
	dbAssociate.Append("SubModels", m, mSub)
	// mSub = &SubModel{Perms: ":::c"}
	// dbAssociate.Append("Many2ManyModels", m, mSub)
	dbAssociate.Close()
	if dbAssociate.Error() != nil {
		t.Error(dbAssociate.Error())
	}

	m = &TestModel{}
	db.Preload("SubModels").First(m)
	db.Close()

	exp := "1"
	got := fmt.Sprint(len(m.SubModels))
	if got != exp {
		t.Fatalf(`expected "%s", got "%s"`, exp, got)
	}

	exp = "1"
	got = fmt.Sprint(m.SubModels[0].OwnerID)
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}

	t.Error()
}

func TestAssociateUniqueCodeGeneration(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables(db)
	setupDBTable(&TestModel{}, db)
	setupDBTable(&SubModel{}, db)

	m := &TestModel{UC: "abc", Field: "a", Perms: ":::"}
	db.Create(m)
	db.Close()

	storage := gormdb.New(dbt, uri)
	storage.UniqueCodeFunc(chassis.RandNumberString)
	storage.UniqueCodeLength(1)

	dbAssociate := storage.AssociateDB(0, []uint{})
	var mSub *SubModel
	for i := 0; i < 15; i++ {
		fmt.Println("\nstep", i)
		mSub = &SubModel{Perms: ":::c"}
		dbAssociate.Append("SubModels", m, mSub)
	}
	dbAssociate.Close()
	if dbAssociate.Error() != nil {
		t.Error(dbAssociate.Error())
	}

	if len(m.UC) <= 1 {
		t.Errorf(`expected "> 1", got "%d"`, len(m.UC))
	}
}
