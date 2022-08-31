package interactor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	ita "github.com/oligoden/chassis/lateral/interactor"
)

func TestSimple(t *testing.T) {
	assert := assert.New(t)

	ita1 := ita.New(0, "a")
	ita2 := ita.New(1, "b")
	ita3 := ita.New(2, "c")

	ita1.Instances = append(ita1.Instances, ita2, ita3)
	ita2.Instances[0] = ita1
	ita2.Instances = append(ita2.Instances, ita3)
	ita3.Instances[0] = ita1
	ita3.Instances[1] = ita2

	ita1.Queue("abc")
	time.Sleep(10 * time.Millisecond)

	r := ita2.Confirmed()
	if assert.NotEmpty(r) {
		assert.Equal("abc", r[0].Query)
	}

	r = ita3.Confirmed()
	if assert.NotEmpty(r) {
		assert.Equal("abc", r[0].Query)
	}
}

func TestConflict(t *testing.T) {
	assert := assert.New(t)

	ita1 := ita.New(0, "a")
	ita2 := ita.New(1, "b")
	ita3 := ita.New(2, "c")

	ita1.Instances = append(ita1.Instances, ita2, ita3)
	ita2.Instances[0] = ita1
	ita2.Instances = append(ita2.Instances, ita3)
	ita3.Instances[0] = ita1
	ita3.Instances[1] = ita2

	ita1.Queue("abc")
	time.Sleep(2 * time.Millisecond)
	ita1.Queue("def")
	time.Sleep(10 * time.Millisecond)

	r := ita1.Confirmed()
	if assert.Len(r, 2) {
		assert.Equal("abc", r[0].Query)
		assert.Equal("def", r[1].Query)
	}

	r = ita2.Confirmed()
	if assert.Len(r, 2) {
		assert.Equal("abc", r[0].Query)
		assert.Equal("def", r[1].Query)
	}

	r = ita3.Confirmed()
	if assert.Len(r, 2) {
		assert.Equal("abc", r[0].Query)
		assert.Equal("def", r[1].Query)
	}
}
