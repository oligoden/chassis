package gormdb

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage"
)

type readDB struct {
	orm              *gorm.DB
	user             uint
	groups           []uint
	err              error
	uniqueCodeLength uint
	uniqueCodeFunc   func(uint, rand.Source) string
	rs               rand.Source
}

func (s Store) ReadDB(user uint, groups []uint) storage.DBReader {
	db := &readDB{}

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

func (db *readDB) Where(s interface{}, d ...interface{}) storage.DBReader {
	db.orm = db.orm.Where(s, d...)
	return db
}

func (db readDB) NewRecord(m interface{}) bool {
	return db.orm.NewRecord(m)
}

type namer interface {
	TableName() string
}

func (db *readDB) First(m interface{}, n ...string) {
	if r, ok := db.read(m, n...); ok {
		r.First(m)
	}
}

func (db *readDB) Last(m interface{}, n ...string) {
	if r, ok := db.read(m, n...); ok {
		r.Last(m)
	}
}

func (db *readDB) Find(m interface{}, n ...string) {
	if r, ok := db.read(m, n...); ok {
		r.Find(m)
	}
}

func (db *readDB) read(m interface{}, n ...string) (*gorm.DB, bool) {
	if db.err != nil {
		return nil, false
	}

	tableName := ""
	if len(n) > 0 {
		tableName = n[0]
	} else {
		mNamer, assertable := m.(namer)
		if !assertable {
			db.err = errors.New("model is not assertable as an table namer")
			return nil, false
		}
		tableName = mNamer.TableName()
	}

	x := db.orm
	db.orm = db.orm.New()

	joins, conditions, selectors := db.readAuthorization(tableName)
	if joins != "" {
		x = x.Joins(joins)
	}
	return x.Where(conditions, selectors...), true
}

func (db readDB) readAuthorization(t string) (string, string, []interface{}) {
	joins := ""
	conditions := fmt.Sprintf("%s.perms LIKE ?", t)
	selectors := []interface{}{"%:%:%:%r%"}

	if db.user != 0 {
		conditions += fmt.Sprintf(" OR %s.perms LIKE ?", t)
		selectors = append(selectors, "%:%:%r%:%")

		joins += fmt.Sprintf("left join groups on groups.owner = %s.owner_id", t)
		joins += fmt.Sprintf(" left join record_groups on record_groups.record_id = %s.hash", t)
		conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND (record_groups.group_id IN (?) OR groups.id IN (?)))", t)
		selectors = append(selectors, "%:%r%:%:%", db.groups, db.groups)

		// conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND owner_id = ?)", t)
		// selectors = append(selectors, "%r%:%:%:%", db.user)

		conditions += fmt.Sprintf(" OR %s.owner_id = ?", t)
		selectors = append(selectors, db.user)
	}

	return joins, conditions, selectors
}

func (db *readDB) Preload(f, t string) storage.DBReader {
	if db.err != nil {
		return db
	}

	db.orm = db.orm.Preload(f, func(pdb *gorm.DB) *gorm.DB {

		joins, conditions, selectors := db.readAuthorization(t)
		if joins != "" {
			pdb = pdb.Joins(joins)
		}
		return pdb.Where(conditions, selectors...)
	})

	return db
}

func (db readDB) ReaderToCreater() storage.DBCreater {
	return &createDB{
		orm:              db.orm,
		user:             db.user,
		groups:           db.groups,
		err:              db.err,
		uniqueCodeLength: db.uniqueCodeLength,
		uniqueCodeFunc:   db.uniqueCodeFunc,
		rs:               db.rs,
	}
}

func (db *readDB) ReaderToUpdater() storage.DBUpdater {
	return &updateDB{
		orm:              db.orm,
		user:             db.user,
		groups:           db.groups,
		err:              db.err,
		uniqueCodeLength: db.uniqueCodeLength,
		uniqueCodeFunc:   db.uniqueCodeFunc,
		rs:               db.rs,
	}
}

func (db *readDB) ReaderToAssociator() storage.DBAssociator {
	return &associateDB{
		orm:              db.orm,
		user:             db.user,
		groups:           db.groups,
		err:              db.err,
		uniqueCodeLength: db.uniqueCodeLength,
		uniqueCodeFunc:   db.uniqueCodeFunc,
		rs:               db.rs,
	}
}

func (db *readDB) Error() error {
	return db.err
}

func (db *readDB) Close() error {
	if db.orm != nil {
		return db.orm.Close()
	}
	return nil
}
