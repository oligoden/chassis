package gormdb

import (
	"errors"

	"github.com/jinzhu/gorm"

	"github.com/oligoden/chassis/storage"
)

type manageDB struct {
	err error
	orm *gorm.DB
}

func (s Store) ManageDB() storage.DBManager {
	db := &manageDB{}

	if s.err != nil {
		db.err = errors.New("store not started")
		return db
	}

	var err error
	db.orm, err = gorm.Open(s.dbt, s.uri)
	if err != nil {
		db.err = err
		return db
	}
	db.orm.LogMode(true)

	return db
}

func (db *manageDB) Manage(m interface{}, action string) {
	if db.err != nil {
		return
	}

	switch action {
	case "migrate":
		db.err = db.orm.AutoMigrate(m).Error
	case "drop":
		db.err = db.orm.DropTable(m).Error
	case "dropIfExists":
		db.err = db.orm.DropTableIfExists(m).Error
	}
}

func (db *manageDB) Error() error {
	return db.err
}

func (db *manageDB) Close() error {
	if db.orm != nil {
		return db.orm.Close()
	}
	return nil
}
