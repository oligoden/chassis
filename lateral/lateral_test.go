package lateral_test

import (
	"testing"
	"teststore"

	"github.com/oligoden/chassis/device/model/data"
	"github.com/oligoden/chassis/lateral"
	"github.com/stretchr/testify/assert"
)

func TestComm(t *testing.T) {
	// ita := interactor.New()
	ltl := lateral.New()
	ltl.SetCommMechanism()

	testStore := teststore.ReadWrite{}
	e := data.Default{}
	ltl.Insert(testStore, e)

	assert.Equal(t, "a", testStore.Data)
}
