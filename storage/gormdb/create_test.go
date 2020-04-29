package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis"
	"github.com/oligoden/chassis/storage/gormdb"
)

func TestUniqueCodeGeneration(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

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
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	// simulate error
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
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	storage := gormdb.New(dbt, uri)
	dbCreate := storage.CreateDB(0, []uint{})

	m := &TestModel{}
	m.Perms = ":::"
	dbCreate.Create(m)
	dbCreate.Close()

	if dbCreate.Error() == nil {
		t.Error(`expected error`)
	}
	exp := "create authorization failed"
	got := dbCreate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreateAuthError(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

	storage := gormdb.New(dbt, uri)
	dbCreate := storage.CreateDB(0, []uint{})

	m := &TestModel{}
	m.Perms = "::"
	dbCreate.Create(m)
	dbCreate.Close()

	if dbCreate.Error() == nil {
		t.Error(`expected error`)
	}
	exp := "the model has incorrect permissions format"
	got := dbCreate.Error().Error()
	if exp != got {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

func TestCreateToUpdate(t *testing.T) {
	cleanDBUserTables()
	setupDBTable(&TestModel{})

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

	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	m = &TestModel{}
	db.First(m)

	if err := db.Close(); err != nil {
		t.Error(err)
	}

	exp := "b"
	got := m.Field
	if got != exp {
		t.Errorf(`expected "%s", got "%s"`, exp, got)
	}
}

// func TestCreateZWithoutOwnGroup(t *testing.T) {
// 	setupDBTable(&Match{})

// 	m := &Match{}
// 	// m.GroupIDs = append(m.GroupIDs, 22)
// 	m.Perms = ":::c"

// 	s := gormdb.NewStore(dbt, uri)

// 	db, err := gorm.Open(dbt, uri)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	db.LogMode(true)
// 	db.Create(&gormdb.User{Username: "public"})
// 	db.Create(&gormdb.Group{Name: "public", Owner: 1})
// 	db.Create(&gormdb.User{Username: "jack"})
// 	db.Create(&gormdb.Group{Name: "jack", Owner: 2})
// 	db.Close()
// 	err = db.Error
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	dbCreate := s.CreateDB(2, []uint{})
// 	dbCreate.Create(m)
// 	dbCreate.Close()
// 	if dbCreate.Error() != nil {
// 		t.Error(dbCreate.Error())
// 	}

// 	if m.MatchID != 1 {
// 		t.Errorf(`expected "1", got "%d"`, m.MatchID)
// 	}
// 	if m.OwnerID != 2 {
// 		t.Errorf(`expected "2", got "%d"`, m.OwnerID)
// 	}
// 	// if !include(m.GroupIDs, 2) {
// 	// 	t.Errorf(`expected "[22 2]", got "%d"`, m.GroupIDs)
// 	// }
// }

// func include(a []uint, x uint) bool {
// 	i := sort.Search(len(a), func(i int) bool { return a[i] == x })
// 	return i < len(a) && a[i] == x
// }

// func TestCreateErrors(t *testing.T) {
// 	s := gormdb.NewStore(dbt, uri)
// 	dbCreate := s.CreateDB(0, []uint{})
// 	dbCreate.Create(&BadModel{})
// 	dbCreate.Close()
// 	if dbCreate.Error() == nil {
// 		t.Error("expected error")
// 	}

// 	db, err := gorm.Open(dbt, uri)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	db.LogMode(true)
// 	db.DropTableIfExists(&Match{})
// 	db.Close()
// 	err = db.Error
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	s = gormdb.NewStore(dbt, uri)
// 	dbCreate = s.CreateDB(0, []uint{})
// 	m := &Match{Field: "test"}
// 	dbCreate.Create(m)
// 	dbCreate.Close()
// 	if dbCreate.Error() == nil {
// 		t.Error("expected error")
// 	}

// 	setupDBTable(&Match{})
// 	s = gormdb.NewStore(dbt, "")
// 	dbCreate = s.CreateDB(0, []uint{})
// 	dbCreate.Create(m)
// 	dbCreate.Close()
// 	if dbCreate.Error() == nil {
// 		t.Error("expected error")
// 	}
// }

// type Match struct {
// 	MatchID  uint     `gorm:"primary_key"`
// 	Field    string   `form:"field"`
// 	Players  []Player `form:"-" json:"players" gorm:"foreignkey:MatchID;association_foreignkey:MatchID"`
// 	UC       string   `gorm:"unique"`
// 	OwnerID  uint
// 	groupIDs []uint
// 	Perms    string
// 	Hash     string
// }

// type Matches []Match

// func (Match) TableName() string {
// 	return "matches"
// }

// func (Matches) TableName() string {
// 	return "matches"
// }

// func (m *Match) UniqueCode(uc ...string) string {
// 	fmt.Println("code", uc)
// 	if len(uc) > 0 {
// 		m.UC = uc[0]
// 	}
// 	return m.UC
// }

// func (m Match) Permissions(p ...string) string {
// 	return m.Perms
// }

// func (m *Match) Owner(o ...uint) uint {
// 	if len(o) > 0 {
// 		m.OwnerID = o[0]
// 	}
// 	return m.OwnerID
// }

// func (m *Match) Groups(g ...uint) []uint {
// 	m.groupIDs = append(m.groupIDs, g...)
// 	return m.groupIDs
// }

// type Player struct {
// 	PlayerID uint `gorm:"primary_key"`
// 	MatchID  uint
// 	Name     string `form:"name"`
// 	UC       string `gorm:"unique"`
// 	OwnerID  uint
// 	groupIDs []uint
// 	Perms    string
// 	Hash     string
// }

// func (Player) TableName() string {
// 	return "players"
// }

// func (m *Player) UniqueCode(uc ...string) string {
// 	if len(uc) > 0 {
// 		m.UC = uc[0]
// 	}
// 	return m.UC
// }

// func (m Player) Permissions(p ...string) string {
// 	return m.Perms
// }

// func (m *Player) Owner(o ...uint) uint {
// 	if len(o) > 0 {
// 		m.OwnerID = o[0]
// 	}
// 	return m.OwnerID
// }

// func (m *Player) Groups(g ...uint) []uint {
// 	m.groupIDs = append(m.groupIDs, g...)
// 	return m.groupIDs
// }
