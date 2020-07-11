package gormdb_test

import (
	"testing"

	"github.com/oligoden/chassis/storage/gormdb"
)

func TestAuthorization(t *testing.T) {
	testCases := []struct {
		desc      string
		user      uint
		groups    []uint
		eGroups   []uint
		eUsers    []uint
		ePerms    string
		operation string
		expAuth   bool
		expErr    string
	}{
		{
			desc:   "too few params",
			ePerms: "::",
			expErr: "the model has incorrect permissions format",
		},
		{
			desc:      "Z",
			ePerms:    ":::u",
			operation: "u",
			expAuth:   true,
		},
		{
			desc:      "Z, must fail, no operation",
			ePerms:    ":::",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "A",
			user:      2,
			ePerms:    "::u:",
			operation: "u",
			expAuth:   true,
		},
		{
			desc:      "A, must fail, unknown user",
			user:      0,
			ePerms:    "::u:",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "G",
			user:      2,
			groups:    []uint{1, 2, 3},
			eGroups:   []uint{2, 4, 6},
			ePerms:    ":u::",
			operation: "u",
			expAuth:   true,
		},
		{
			desc:      "G, must fail, unknown user",
			user:      0,
			groups:    []uint{1},
			eGroups:   []uint{1},
			ePerms:    ":u::",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "G, must fail, wrong group",
			user:      2,
			groups:    []uint{2},
			eGroups:   []uint{1},
			ePerms:    ":u::",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "U",
			user:      2,
			eUsers:    []uint{1, 2, 3},
			ePerms:    "u:::",
			operation: "u",
			expAuth:   true,
		},
		{
			desc:      "U, must fail, unknown user",
			user:      0,
			eUsers:    []uint{1, 2, 3},
			ePerms:    "u:::",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "U, must fail, wrong user",
			user:      4,
			eUsers:    []uint{1, 2, 3},
			ePerms:    "u:::",
			operation: "u",
			expAuth:   false,
		},
		{
			desc:      "O",
			user:      1,
			ePerms:    ":::",
			operation: "u",
			expAuth:   true,
		},
		{
			desc:      "O",
			user:      1,
			ePerms:    ":::",
			operation: "c",
			expAuth:   false,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			e := &TestModel{}
			e.Perms = tC.ePerms
			e.Users(tC.eUsers...)
			e.Groups(tC.eGroups...)
			e.OwnerID = 1

			auth, err := gormdb.Authorize(e, tC.operation, tC.user, tC.groups)
			if tC.expErr != "" && err == nil {
				t.Errorf(`expected error "%s"`, tC.expErr)
			} else if tC.expErr == "" && err != nil {
				t.Errorf(`expected no error, got "%s"`, err)
			} else if err != nil && err.Error() != tC.expErr {
				t.Errorf(`expected error "%s", got "%s"`, tC.expErr, err)
			}
			if auth != tC.expAuth {
				t.Errorf(`expected auth "%t", got "%t"`, tC.expAuth, auth)
			}
		})
	}
}
