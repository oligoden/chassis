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
	// Error() error
}

type dbReader interface {
	Where(interface{}, ...interface{}) DBReader
	First(interface{}, ...string)
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
	// Error() error
}

type dbAssociator interface {
	Append(string, Authenticator, Authenticator)
	// Clear(Authenticator, string)
}
type DBAssociator interface {
	dbManager
	dbAssociator
	// Error() error
}

type DBManager interface {
	dbManager
	Manage(interface{}, string)
	// Error() error
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
}
