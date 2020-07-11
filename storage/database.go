package storage

type Storer interface {
	ManageDB() DBManager
	CreateDB(uint, []uint) DBCreater
	ReadDB(uint, []uint) DBReader
	// UpdateDB(uint, []uint) DBUpdater
	// AssociationDB(uint, []uint) DBAssociater
	UniqueCodeLength(...uint) uint
	Error() error
}

type DBCreater interface {
	dbManager
	Create(Authenticator)
	CreaterToUpdater() DBUpdater
}

type DBReader interface {
	dbManager
	dbReader
}

type dbReader interface {
	Where(interface{}, ...interface{}) DBReader
	First(interface{}, ...string)
	Last(interface{}, ...string)
	Find(interface{}, ...string)
	Preload(string, string) DBReader
	NewRecord(interface{}) bool
	ReaderToAssociator() DBAssociator
	ReaderToCreater() DBCreater
	ReaderToUpdater() DBUpdater
}

type dbUpdater interface {
	Save(Authenticator, ...string)
}

type DBUpdater interface {
	dbManager
	dbUpdater
}

type dbAssociator interface {
	Append(string, Authenticator, Authenticator)
	// Clear(Authenticator, string)
}
type DBAssociator interface {
	dbManager
	dbAssociator
}

type DBManager interface {
	dbManager
	Manage(interface{}, string)
}

type dbManager interface {
	// Open() *DB
	Close() error
	// DropTable() *DB
	Error() error
}

type Authenticator interface {
	UniqueCode(...string) string
	Permissions(...string) string
	Owner(...uint) uint
	Groups(...uint) []uint
	Users(...uint) []uint
}
