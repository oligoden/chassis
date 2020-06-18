package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestReadAuthorization(t *testing.T) {
	t.SkipNow()
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables()
	storage := gormdb.New(dbt, uri)

	mGroup := &gormdb.Group{Owner: 1}
	db.Create(mGroup)
	mGroup = &gormdb.Group{Owner: 2}
	db.Create(mGroup)
	mRecordGroup := &gormdb.RecordGroup{
		GroupID:  1,
		RecordID: "a",
		Owner:    2,
	}
	db.Create(mRecordGroup)
	setupDBTable(&TestModel{}, db)

	testCases := []struct {
		desc       string
		user       uint
		recOwnerID uint
		groups     []uint
		perms      string
		setField   string
		expField   string
	}{
		{
			desc:     "Pass_Z",
			user:     0,
			groups:   []uint{},
			perms:    ":::r",
			setField: "a",
			expField: "a",
		},
		{
			desc:     "Fail_Z",
			user:     0,
			groups:   []uint{},
			perms:    ":::",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Pass_A",
			user:     1,
			groups:   []uint{},
			perms:    "::r:",
			setField: "b",
			expField: "b",
		},
		{
			desc:     "Fail1_A_missing_permission",
			user:     1,
			groups:   []uint{},
			perms:    ":::",
			setField: "b",
			expField: "",
		},
		{
			desc:     "Fail2_A_missing_userID",
			user:     0,
			groups:   []uint{},
			perms:    "::r:",
			setField: "b",
			expField: "",
		},
		{
			desc:       "Pass_G",
			user:       1,
			groups:     []uint{2},
			perms:      ":r::",
			recOwnerID: 2,
			setField:   "a",
			expField:   "a",
		},
		{
			desc:       "Pass_G_RecordGroup",
			user:       1,
			groups:     []uint{1},
			perms:      ":r::",
			recOwnerID: 2,
			setField:   "a",
			expField:   "a",
		},
		{
			desc:     "Fail_G_missing_group",
			user:     1,
			groups:   []uint{},
			perms:    ":r::",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Fail_G_missing_permission",
			user:     1,
			groups:   []uint{2},
			perms:    ":::",
			setField: "a",
			expField: "",
		},
		{
			desc:       "Pass_O",
			user:       1,
			groups:     []uint{},
			perms:      ":::",
			recOwnerID: 1,
			setField:   "a",
			expField:   "a",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			setupDBTable(&TestModel{}, db)
			m := &TestModel{
				Field:   tC.expField,
				Hash:    tC.expField,
				Perms:   tC.perms,
				OwnerID: tC.recOwnerID,
			}
			db.Create(m)

			dbRead := storage.ReadDB(tC.user, tC.groups)
			m = &TestModel{}
			dbRead.First(m)
			dbRead.Close()
			if dbRead.Error() != nil {
				t.Error(dbRead.Error())
			}

			exp := tC.expField
			got := m.Field
			if got != exp {
				t.Errorf(`expected "%s", got "%s"`, exp, got)
			}
		})
	}

	db.Close()
}
