package gormdb_test

import (
	"testing"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestAuthorizationFewPerms(t *testing.T) {
	m := &TestModel{}
	m.Perms = "::"

	auth, err := gormdb.Authorize(m, "", 0, []uint{})
	if !(err != nil && err.Error() == "the model has incorrect permissions format") {
		t.Error(`expected incorrect permissions format error, got`, err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationONegOwner(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = "c:::"

	auth, err := gormdb.Authorize(m, "c", 0, []uint{})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationONegPerms(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationONegUser(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = "c:::"

	auth, err := gormdb.Authorize(m, "c", 2, []uint{})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationO(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = "c:::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{})
	if err != nil {
		t.Error(err)
	}
	if !auth {
		t.Error(`expected true, got false`)
	}
}

func TestAuthorizationGNegOwner(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":c::"

	auth, err := gormdb.Authorize(m, "c", 0, []uint{1})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationGNegPerms(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{1})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationGNegGroup(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":c::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{1, 3})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationG(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.groupIDs = []uint{7}
	m.Perms = ":c::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{5, 1, 7})
	if err != nil {
		t.Error(err)
	}
	if !auth {
		t.Error(`expected true, got false`)
	}
}

func TestAuthorizationANegOwner(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = "::c:"

	auth, err := gormdb.Authorize(m, "c", 0, []uint{1})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationANegPerms(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{1})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationA(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = "::c:"

	auth, err := gormdb.Authorize(m, "c", 2, []uint{5, 1, 7})
	if err != nil {
		t.Error(err)
	}
	if !auth {
		t.Error(`expected true, got false`)
	}
}

func TestAuthorizationZNegPerms(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::"

	auth, err := gormdb.Authorize(m, "c", 1, []uint{1})
	if err != nil {
		t.Error(err)
	}
	if auth {
		t.Error(`expected false, got true`)
	}
}

func TestAuthorizationZ(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::c"

	auth, err := gormdb.Authorize(m, "c", 0, []uint{5, 1, 7})
	if err != nil {
		t.Error(err)
	}
	if !auth {
		t.Error(`expected true, got false`)
	}
}

func TestAuthorizationZUser(t *testing.T) {
	m := &TestModel{}
	m.OwnerID = 1
	m.Perms = ":::c"

	auth, err := gormdb.Authorize(m, "c", 2, []uint{5, 1, 7})
	if err != nil {
		t.Error(err)
	}
	if !auth {
		t.Error(`expected true, got false`)
	}
}

// func TestAuthorizationStoreNegPerms(t *testing.T) {
// 	setupDBTable(&Match{})

// 	recMatch := &Match{}
// 	recMatch.perms = "::"

// 	s := gormdb.NewStore(dbt, uri)
// 	dbCreate := s.CreateDB(0, []uint{})
// 	err := dbCreate.Create(recMatch).Error()

// 	if !(err != nil && err.Error() == "the model has less than four permissions") {
// 		t.Error(`expected less than four permissions error, got`, err)
// 	}
// 	dbCreate.Close()
// }

// func TestAuthorizationZNeg(t *testing.T) {
// 	setupDBTable(&Match{})
// 	store := store.NewStore(dbt, uri)

// 	dbCreate := store.CreateDB(0, []uint{})
// 	recMatch := &Match{}
// 	recMatch.Owner = "jack"
// 	recMatch.Group = "jack"
// 	recMatch.Perms = ":::r"
// 	err := dbCreate.Create(recMatch).Error()
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	dbRead := store.ReadDB(0, []uint{})
// 	recMatch = &Match{}
// 	err = dbRead.First(recMatch).Error()
// 	if !(err != nil && err.Error() == "record not found") {
// 		t.Error(`expected no record, error:`, err)
// 	}
// 	dbRead.Close()

// 	dropDBTable(&Match{})
// }
