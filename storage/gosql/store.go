package gosql

import (
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oligoden/chassis"
)

type Store struct {
	dbt              string
	uri              string
	err              error
	uniqueCodeLength uint
	ucFunc           func(uint) string
	rnd              *rand.Rand
}

func New(uri string) *Store {
	s := new(Store)

	db, err := sql.Open("mysql", uri)
	if err != nil {
		s.err = fmt.Errorf("opening db connection for new store migration: %w", err)
		return s
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	s.dbt = "mysql"
	s.uri = uri

	s.uniqueCodeLength = 2
	rs := rand.NewSource(time.Now().UnixNano())
	s.rnd = rand.New(rs)
	s.ucFunc = s.randString

	_, err = db.Exec("CREATE TABLE `users` (`owner_id` int unsigned AUTO_INCREMENT,`uc` varchar(255) UNIQUE NOT NULL DEFAULT '',`ts` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,`username` varchar(255) NOT NULL DEFAULT '',`pass_hash` varchar(255) NOT NULL DEFAULT '',`salt` varchar(255) NOT NULL DEFAULT '',`perms` varchar(255),`hash` varchar(255) NOT NULL DEFAULT '', PRIMARY KEY (`owner_id`))")
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1050") {
			s.err = fmt.Errorf("doing new store db migration: %w", err)
			return s
		}
	}

	_, err = db.Exec("CREATE TABLE `groups` (`id` int unsigned AUTO_INCREMENT,`ts` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,`name` varchar(255),`owner` int unsigned,`perms` varchar(255) , PRIMARY KEY (`id`))")
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1050") {
			s.err = fmt.Errorf("doing new store db migration: %w", err)
			return s
		}
	}

	_, err = db.Exec("CREATE TABLE `record_groups` (`id` int unsigned AUTO_INCREMENT,`ts` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,`record_id` varchar(255),`group_id` int unsigned,`owner` int unsigned,`perms` varchar(255) , PRIMARY KEY (`id`))")
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1050") {
			s.err = fmt.Errorf("doing new store db migration: %w", err)
			return s
		}
	}

	_, err = db.Exec("CREATE TABLE `record_users` (`id` int unsigned AUTO_INCREMENT,`ts` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,`description` varchar(255),`record_id` varchar(255),`user_id` int unsigned,`owner` int unsigned,`perms` varchar(255) , PRIMARY KEY (`id`))")
	if err != nil {
		if !strings.Contains(err.Error(), "Error 1050") {
			s.err = fmt.Errorf("doing new store db migration: %w", err)
			return s
		}
	}

	return s
}

func ConnURL(serviceName ...string) string {
	pre := ""

	if len(serviceName) >= 1 {
		pre = strings.ToUpper(serviceName[0])
		if !strings.HasSuffix(pre, "_") {
			pre = pre + "_"
		}
	}

	dbUser := os.Getenv(pre + "DB_USER")
	dbPass := os.Getenv(pre + "DB_PASS")
	dbAddr := os.Getenv(pre + "DB_ADDR")
	dbPort := os.Getenv(pre + "DB_PORT")
	dbName := os.Getenv(pre + "DB_NAME")

	if dbAddr == "" {
		dbAddr = "localhost"
	}

	if dbPort == "" {
		dbPort = "3306"
	}

	params := "charset=utf8&parseTime=True&loc=Local"
	format := "%s:%s@tcp(%s:%s)/%s?%s"

	return fmt.Sprintf(format, dbUser, dbPass, dbAddr, dbPort, dbName, params)
}

type Migrater interface {
	Migrate(*sql.DB) error
}

func (s *Store) Migrate(e Migrater) {
	if s.err != nil {
		return
	}

	db, err := sql.Open(s.dbt, s.uri)
	if err != nil {
		s.err = chassis.Mark("opening db connection for migration", err)
	}
	defer db.Close()

	err = e.Migrate(db)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return
		}
		s.err = chassis.Mark("running db migration", err)
	}
}

func (s *Store) UniqueCodeLength(ucl ...uint) uint {
	if len(ucl) > 0 {
		s.uniqueCodeLength = ucl[0]
	}
	return s.uniqueCodeLength
}

func (s *Store) UniqueCodeFunc(ucf ...func(uint) string) func(uint) string {
	if len(ucf) > 0 {
		s.ucFunc = ucf[0]
	}
	return s.ucFunc
}

const (
	numalphaLetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits       = 6
	letterIdxMask       = 1<<letterIdxBits - 1
	letterIdxMax        = 63 / letterIdxBits
)

func (s Store) randString(n uint) string {
	return randString(n, s.rnd)
}

// randString generates a random string
func randString(n uint, rnd *rand.Rand) string {
	// solution from http://stackoverflow.com/a/31832326
	b := make([]byte, n)
	for i, cache, remain := int(n-1), rnd.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rnd.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(numalphaLetterBytes) {
			b[i] = numalphaLetterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func (s Store) Rnd() *rand.Rand {
	return s.rnd
}

func (s Store) Err() error {
	return s.err
}

type User struct {
	OwnerID  uint      `gosql:"primary_key,read-only" json:"-"`
	UC       string    `gosql:"unique" json:"uc" form:"uc"`
	TS       time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
	Username string    `gosql:"not null" json:"username"`
	Password string    `gosql:"-" json:"-" form:"password"`
	PassHash string    `json:"-"`
	Salt     string    `json:"salt"`
	// 	UserGroups []Group   `gosql:"many2many:user_groups"`
	GroupIDs []uint `gosql:"-" json:"-"`
	UserIDs  []uint `gosql:"-" json:"-"`
	Perms    string `json:"-"`
	Hash     string `json:"-"`
	rnd      *rand.Rand
	req      *http.Request
}

func NewUserRecord(req *http.Request, rnd *rand.Rand) *User {
	r := &User{}
	r.req = req
	r.rnd = rnd
	r.Perms = ":::cr"
	return r
}

func (User) TableName() string {
	return "users"
}

func (e *User) Prepare() error {
	h := sha256.New()
	h.Write([]byte(e.Password + e.Salt))
	bs := h.Sum(nil)
	e.PassHash = fmt.Sprintf("%x", bs)

	e.Username = randString(6, e.rnd)
	return nil
}

func (e User) Complete() error {
	e.req.Header.Set("X_user", fmt.Sprint(e.IDValue()))
	return nil
}

func (e *User) IDValue(id ...uint) uint {
	if len(id) > 0 {
		e.OwnerID = id[0]
	}
	return e.OwnerID
}

func (e *User) UniqueCode(uc ...string) string {
	if len(uc) > 0 {
		e.UC = uc[0]
	}
	return e.UC
}

func (e *User) Permissions(p ...string) string {
	if len(p) > 0 {
		e.Perms = p[0]
	}
	return e.Perms
}

func (e *User) Owner(o ...uint) uint {
	if len(o) > 0 {
		e.OwnerID = o[0]
	}
	return e.OwnerID
}

func (e *User) Groups(g ...uint) []uint {
	e.GroupIDs = append(e.GroupIDs, g...)
	return e.GroupIDs
}

func (e *User) Users(u ...uint) []uint {
	e.UserIDs = append(e.UserIDs, u...)
	return e.UserIDs
}

func (e *User) Hasher() error {
	json, err := json.Marshal(e)
	if err != nil {
		return err
	}
	h := sha1.New()
	h.Write(json)
	e.Hash = fmt.Sprintf("%x", h.Sum(nil))

	return nil
}

type UserRecords []User

func (UserRecords) TableName() string {
	return "users"
}

func (e UserRecords) Permissions(p ...string) string {
	return ""
}

func (e UserRecords) Owner(o ...uint) uint {
	return 0
}

func (e UserRecords) Users(u ...uint) []uint {
	return []uint{}
}

func (e UserRecords) Groups(g ...uint) []uint {
	return []uint{}
}

func (UserRecords) IDValue(...uint) uint {
	return 0
}

func (e UserRecords) UniqueCode(uc ...string) string {
	return ""
}

// type Group struct {
// 	ID    uint      `gorm:"primary_key"`
// 	TS    time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
// 	Name  string
// 	Owner uint
// 	Perms string
// }

// func (Group) TableName() string {
// 	return "groups"
// }

// type RecordGroup struct {
// 	ID       uint      `gorm:"primary_key"`
// 	TS       time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
// 	RecordID string
// 	GroupID  uint
// 	Owner    uint
// 	Perms    string
// }

// type RecordUser struct {
// 	ID          uint      `gorm:"primary_key"`
// 	TS          time.Time `sql:"DEFAULT:CURRENT_TIMESTAMP"`
// 	Description string
// 	RecordID    string
// 	UserID      uint
// 	Owner       uint
// 	Perms       string
// }
