package data

import "github.com/oligoden/chassis/storage"

type List struct {
}

func (List) Prepare() error {
	return nil
}

func (e *List) Read(db storage.DBReader) error {
	return nil
}

func (e List) TableName() string {
	return "models"
}

func (e *List) Complete() error {
	return nil
}

func (e List) Hasher() error {
	return nil
}

func (e List) UniqueCode(uc ...string) string {
	return ""
}

func (e List) Permissions(p ...string) string {
	return ""
}

func (e List) Owner(o ...uint) uint {
	return 0
}

func (e List) Users(us ...uint) []uint {
	return []uint{}
}

func (e List) Groups(gs ...uint) []uint {
	return []uint{}
}
