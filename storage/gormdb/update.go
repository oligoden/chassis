package gormdb

import (
	"errors"
	"math/rand"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage"
)

type updateDB struct {
	orm              *gorm.DB
	user             uint
	groups           []uint
	err              error
	uniqueCodeLength uint
	rs               rand.Source
}

func (s Store) UpdateDB(user uint, groups []uint) storage.DBUpdater {
	db := &updateDB{}

	if s.err != nil {
		db.err = errors.New("store error")
		return db
	}

	db.orm, db.err = gorm.Open(s.dbt, s.uri)
	if db.err != nil {
		return nil
	}
	db.orm.LogMode(true)
	db.user = user
	db.groups = groups

	db.uniqueCodeLength = s.UniqueCodeLength()
	db.rs = rand.New(s.rs)

	return db
}

func (db *updateDB) Save(m storage.Authenticator, params ...string) {
	if db.err != nil {
		return
	}

	perm := "u"

	for _, param := range params {
		if param == "with-create" {
			perm = "c"
		}
	}

	if auth, err := Authorize(m, perm, db.user, db.groups); !auth {
		if err != nil {
			db.err = err
			return
		}
		db.err = errors.New("update authorization failed")
		return
	}

	db.orm = db.orm.Save(m)
}

func (db *updateDB) Error() error {
	return db.err
}

func (db *updateDB) Close() error {
	if db.orm != nil {
		return db.orm.Close()
	}
	return nil
}
