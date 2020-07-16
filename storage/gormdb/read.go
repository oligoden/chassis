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

func (db *readDB) First(m interface{}, params ...string) {
	if r, ok := db.read(m, params...); ok {
		r.First(m)
	}
}

func (db *readDB) Last(m interface{}, params ...string) {
	if r, ok := db.read(m, params...); ok {
		r.Last(m)
	}
}

func (db *readDB) Find(m interface{}, params ...string) {
	if r, ok := db.read(m, params...); ok {
		r.Find(m)
	}
}

func (db *readDB) read(m interface{}, params ...string) (*gorm.DB, bool) {
	if db.err != nil {
		return nil, false
	}

	authParams := []string{}
	tableName := ""
	for _, param := range params {
		if param == "with-update" {
			authParams = append(authParams, param)
		} else {
			tableName = param
		}
	}

	if tableName == "" {
		mNamer, assertable := m.(namer)
		if !assertable {
			db.err = errors.New("model is not assertable as an table namer")
			return nil, false
		}
		tableName = mNamer.TableName()
	}

	x := db.orm
	db.orm = db.orm.New()

	joins, conditions, selectors := db.readAuthorization(tableName, params...)
	if joins != "" {
		x = x.Joins(joins)
	}
	return x.Where(conditions, selectors...), true
}

func (db readDB) readAuthorization(t string, params ...string) (string, string, []interface{}) {
	perm := "r"

	for _, param := range params {
		if param == "with-update" {
			perm = "u"
		}
	}

	permsZ := fmt.Sprintf("%%:%%:%%:%%%s%%", perm)
	permsA := fmt.Sprintf("%%:%%:%%%s%%:%%", perm)
	permsG := fmt.Sprintf("%%:%%%s%%:%%:%%", perm)
	permsU := fmt.Sprintf("%%%s%%:%%:%%:%%", perm)

	joins := ""
	conditions := fmt.Sprintf("%s.perms LIKE ?", t)
	selectors := []interface{}{permsZ}

	if db.user != 0 {
		conditions += fmt.Sprintf(" OR %s.perms LIKE ?", t)
		selectors = append(selectors, permsA)

		joins += fmt.Sprintf(" left join record_groups on record_groups.record_id = %s.hash", t)
		conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND record_groups.group_id IN (?))", t)
		selectors = append(selectors, permsG, db.groups)

		joins += fmt.Sprintf(" left join record_users on record_users.record_id = %s.hash", t)
		conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND record_users.user_id = ?)", t)
		selectors = append(selectors, permsU, db.user)

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
