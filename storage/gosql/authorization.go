package gosql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/oligoden/chassis/storage"
)

func Authorize(e storage.Authenticator, p string, user uint, groups []uint) (bool, string, error) {
	perms := strings.Split(e.Permissions(), ":")
	if len(perms) != 4 {
		return false, "", errors.New("the model has incorrect permissions format")
	}

	if e.Owner() == user && p != "c" {
		return true, "", nil
	}

	if strings.Contains(perms[3], p) {
		return true, "", nil
	}

	if user != 0 {
		if strings.Contains(perms[2], p) {
			return true, "", nil
		} else if strings.Contains(perms[1], p) {
			for _, g := range e.Groups() {
				for _, gi := range groups {
					if g == gi {
						return true, "", nil
					}
				}
			}
		} else if strings.Contains(perms[0], p) {
			for _, u := range e.Users() {
				if u == user {
					return true, "", nil
				}
			}
		}
	}

	if e.Owner() != user {
		return false, fmt.Sprintf("not owner of record (owner:%d, user:%d)", e.Owner(), user), nil
	}

	return false, "", nil
}
