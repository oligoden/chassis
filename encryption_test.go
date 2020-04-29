package chassis_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/oligoden/chassis"
)

func Test(t *testing.T) {
	rs := rand.NewSource(time.Now().UnixNano())

	if len(chassis.RandString(4, rand.New(rs))) != 4 {
		t.Error("incorrect length")
	}
}
