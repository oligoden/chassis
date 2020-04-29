package gormdb

import (
	"errors"
	"math/rand"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage"
)

type createDB struct {
	orm              *gorm.DB
	user             uint
	groups           []uint
	err              error
	uniqueCodeLength uint
	uniqueCodeFunc   func(uint, rand.Source) string
	rs               rand.Source
}

func (s Store) CreateDB(user uint, groups []uint) storage.DBCreater {
	db := &createDB{}

	if s.err != nil {
		db.err = errors.New("store error")
		return db
	}

	db.orm, db.err = gorm.Open(s.dbt, s.uri)
	if db.err != nil {
		return db
	}
	db.orm.LogMode(true)
	db.user = user
	db.groups = groups

	db.uniqueCodeLength = s.UniqueCodeLength()
	db.uniqueCodeFunc = s.UniqueCodeFunc()
	db.rs = rand.New(s.rs)

	return db
}

func (db *createDB) Create(m storage.Authenticator) {
	if db.err != nil {
		return
	}

	m.Owner(db.user)

	if auth, err := Authorize(m, "c", db.user, db.groups); !auth {
		if err != nil {
			db.err = err
			return
		}
		db.err = errors.New("create authorization failed")
		return
	}

	m.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
	err := db.orm.Create(m).Error
	if err != nil {
		if isDuplicateUniqueCode(err) {
			db.err = db.retryCreate(m)
		} else {
			db.err = err
		}
	}
	return
}

func isDuplicateUniqueCode(err error) bool {
	return strings.Contains(err.Error(), "Error 1062: Duplicate entry") && strings.Contains(err.Error(), "for key 'uc'")
}

func (db *createDB) retryCreate(m storage.Authenticator) error {
	for i := 0; i < 3; i++ {
		m.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
		err := db.orm.Create(m).Error
		if err != nil {
			if !isDuplicateUniqueCode(err) {
				return err
			}
		} else {
			return nil
		}
	}

	db.uniqueCodeLength++
	m.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
	err := db.orm.Create(m).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *createDB) CreaterToUpdater() storage.DBUpdater {
	return &updateDB{
		orm:    db.orm,
		user:   db.user,
		groups: db.groups,
		err:    db.err,
	}
}

func (db *createDB) Error() error {
	return db.err
}

func (db *createDB) Close() error {
	if db.orm != nil {
		return db.orm.Close()
	}
	return nil
}
