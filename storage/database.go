package storage

type Crudder interface {
	// Where(WhereSetter)
	Create(Operator)
	Read(Operator)
	Update(Operator)
	Delete(Operator)
	AddModifiers(...Modifier)
	Err(...error) error
}

type Modifier interface {
	Compile(...string) (string, []interface{})
	Order() int
}

// type Storer interface {
// 	Connect(Identificator) DoCloser
// 	ManageDB() DBManager
// 	CreateDB(uint, []uint) DBCreater
// 	ReadDB(uint, []uint) DBReader
// 	UpdateDB(uint, []uint) DBUpdater
// 	AssociateDB(uint, []uint) DBAssociator
// 	UniqueCodeLength(...uint) uint
// 	Error() error
// }

type Identificator interface {
	User() (uint, []uint)
}

type DoCloser interface {
	Do(string, interface{}, ...string) DoCloser
	Close() error
	Err() error
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
	Delete(string, Authenticator, Authenticator)
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
	IDValue(...uint) uint
	UniqueCode(...string) string
	Permissions(...string) string
	Owner(...uint) uint
	Groups(...uint) []uint
	Users(...uint) []uint
}

type TableNamer interface {
	TableName() string
}

type Operator interface {
	Authenticator
	TableNamer
}
