package gormdb_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/oligoden/chassis/storage/gormdb"
)

func TestReadAuthorization(t *testing.T) {
	db, err := gorm.Open(dbt, uri)
	if err != nil {
		t.Error(err)
	}
	db.LogMode(true)

	cleanDBUserTables()
	storage := gormdb.New(dbt, uri)

	mRecordGroup := &gormdb.RecordGroup{
		GroupID:  1,
		RecordID: "x",
	}
	db.Create(mRecordGroup)
	mRecordUser := &gormdb.RecordUser{
		RecordID: "x",
		UserID:   2,
	}
	db.Create(mRecordUser)

	testCases := []struct {
		desc     string
		user     uint
		groups   []uint
		recPerms string
		recHash  string
		setField string
		expField string
	}{
		{
			desc:     "Pass_Z",
			recPerms: ":::r",
			setField: "a",
			expField: "a",
		},
		{
			desc:     "Fail_Z",
			recPerms: ":::",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Pass_A",
			user:     2,
			recPerms: "::r:",
			setField: "b",
			expField: "b",
		},
		{
			desc:     "Fail_A_missing_permission",
			user:     2,
			recPerms: ":::",
			setField: "b",
			expField: "",
		},
		{
			desc:     "Fail_A_missing_userID",
			user:     0,
			recPerms: "::r:",
			setField: "b",
			expField: "",
		},
		{
			desc:     "Pass_G",
			user:     2,
			groups:   []uint{1},
			recPerms: ":r::",
			recHash:  "x",
			setField: "a",
			expField: "a",
		},
		{
			desc:     "Fail_G_missing_group",
			user:     2,
			groups:   []uint{},
			recPerms: ":r::",
			recHash:  "x",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Fail_G_missing_permission",
			user:     2,
			groups:   []uint{2},
			recPerms: ":::",
			recHash:  "x",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Pass_U",
			user:     2,
			recPerms: "r:::",
			recHash:  "x",
			setField: "a",
			expField: "a",
		},
		{
			desc:     "Fail_U_missing_wrong_user",
			user:     3,
			recPerms: "r:::",
			recHash:  "x",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Fail_U_missing_permission",
			user:     2,
			recPerms: ":::",
			recHash:  "x",
			setField: "a",
			expField: "",
		},
		{
			desc:     "Pass_O",
			user:     1,
			recPerms: ":::",
			setField: "a",
			expField: "a",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			setupDBTable(&TestModel{}, db)

			m := &TestModel{
				Field:   tC.expField,
				Hash:    tC.recHash,
				Perms:   tC.recPerms,
				OwnerID: 1,
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
