package data

import "github.com/oligoden/chassis/storage"

type List struct {
}

func (List) Prepare() error {
	return nil
}

func (m *List) Read(db storage.DBReader) error {
	return nil
}

func (m List) TableName() string {
	return "models"
}

func (m *List) Complete() error {
	return nil
}

func (m List) Hasher() error {
	return nil
}

func (m List) UniqueCode(uc ...string) string {
	return ""
}

func (m List) Permissions(p ...string) string {
	return ""
}

func (m List) Owner(o ...uint) uint {
	return 0
}

func (m List) Groups(gs ...uint) []uint {
	return []uint{}
}
