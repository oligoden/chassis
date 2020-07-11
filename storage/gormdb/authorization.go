package gormdb

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oligoden/chassis/storage"
)

func Authorize(m storage.Authenticator, p string, user uint, groups []uint) (bool, error) {
	perms := strings.Split(m.Permissions(), ":")
	if len(perms) != 4 {
		return false, errors.New("the model has incorrect permissions format")
	}

	if m.Owner() == user && p != "c" {
		fmt.Println("owner", m.Owner(), user)
		return true, nil
	}

	if strings.Contains(perms[3], p) {
		return true, nil
	}

	if user != 0 {
		if strings.Contains(perms[2], p) {
			return true, nil
		} else if strings.Contains(perms[1], p) {
			for _, g := range m.Groups() {
				for _, gi := range groups {
					if g == gi {
						return true, nil
					}
				}
			}
		} else if strings.Contains(perms[0], p) {
			for _, u := range m.Users() {
				if u == user {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
