package gormdb

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"
	"github.com/oligoden/chassis/storage"
)

type associateDB struct {
	orm              *gorm.DB
	assoc            *gorm.Association
	user             uint
	groups           []uint
	err              error
	uniqueCodeLength uint
	uniqueCodeFunc   func(uint, rand.Source) string
	rs               rand.Source
}

func (s *Store) AssociateDB(user uint, groups []uint) storage.DBAssociator {
	db := &associateDB{}

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

func (db *associateDB) Append(f string, m, a storage.Authenticator) {
	if db.err != nil {
		return
	}

	m.Owner(db.user)

	if auth, err := Authorize(a, "c", db.user, db.groups); !auth {
		if err != nil {
			db.err = err
			return
		}
		db.err = errors.New("associate create authorization failed")
		return
	}

	a.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
	err := db.orm.Create(a).Error
	if err != nil {
		if isDuplicateUniqueCode(err) {
			db.err = db.retryAssociate(f, m, a)
		} else {
			db.err = err
		}
	}

	db.orm.Model(m).Association(f).Append(a)
}

func (db *associateDB) retryAssociate(f string, m, a storage.Authenticator) error {
	for i := 0; i < 3; i++ {
		fmt.Println("attempt", i)
		a.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
		err := db.orm.Create(a).Error
		if err != nil {
			if !isDuplicateUniqueCode(err) {
				return err
			}
		} else {
			return nil
		}
	}

	db.uniqueCodeLength++
	fmt.Println("increased to", db.uniqueCodeLength)
	a.UniqueCode(db.uniqueCodeFunc(db.uniqueCodeLength, db.rs))
	fmt.Println("code", a.UniqueCode())
	err := db.orm.Create(a).Error
	if err != nil {
		return err
	}
	return nil
}

func (db *associateDB) Error() error {
	return db.err
}

func (db *associateDB) Close() error {
	if db.orm != nil {
		return db.orm.Close()
	}
	return nil
}
