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

func (db *readDB) First(m interface{}, n ...string) {
	if db.err != nil {
		return
	}

	tableName := "testmodels"

	x := db.orm
	db.orm = db.orm.New()

	conditions := fmt.Sprintf("%s.perms LIKE ?", tableName)
	selectors := []interface{}{"%:%:%:%r%"}

	if db.user != 0 {
		conditions += fmt.Sprintf(" OR %s.perms LIKE ?", tableName)
		selectors = append(selectors, "%:%:%r%:%")

		recordGroupJoin := fmt.Sprintf("left join groups on groups.owner = %s.owner_id", tableName)
		recordGroupJoin += fmt.Sprintf(" left join record_groups on record_groups.record_id = %s.hash", tableName)
		conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND (record_groups.group_id IN (?) OR groups.id IN (?)))", tableName)
		selectors = append(selectors, "%:%r%:%:%", db.groups, db.groups)
		x = x.Joins(recordGroupJoin)

		// conditions += fmt.Sprintf(" OR (%s.perms LIKE ? AND owner_id = ?)", tableName)
		// selectors = append(selectors, "%r%:%:%:%", db.user)

		conditions += fmt.Sprintf(" OR %s.owner_id = ?", tableName)
		selectors = append(selectors, db.user)
	}

	x.Where(conditions, selectors...).First(m)
}

func (db *readDB) Find(interface{}, ...string) {
	return
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
