package bitset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert := assert.New(t)

	bs := New()
	assert.NotNil(bs)

	bs.b = append(bs.b, 0b00000101)
	assert.True(bs.Get(0))
	assert.False(bs.Get(1))
	assert.True(bs.Get(2))

	bs.b = append(bs.b, 0b00000100)
	assert.True(bs.Get(10))

	bs.Set(12, true)
	assert.Equal(byte(0b00010100), bs.b[1])
	assert.Equal(13, bs.Len())

	bs.Set(12, false)
	assert.Equal(byte(0b00000100), bs.b[1])

	bs.Set(16, true)
	assert.Equal(byte(0b00000001), bs.b[2])
	assert.Equal(17, bs.Len())

	assert.Equal(9, bs.Len(9))
	assert.Len(bs.b, 2)
	assert.Equal(byte(0b00000000), bs.b[1])

	assert.Equal(19, bs.Len(19))
	assert.Len(bs.b, 3)
	assert.Equal(byte(0b00000000), bs.b[2])

	assert.False(bs.Zero())
	bs.Set(0, false)
	bs.Set(2, false)
	assert.True(bs.Zero())
}

func TestOr(t *testing.T) {
	assert := assert.New(t)

	bsA := New()
	bsB := New()
	bsC := New()
	bsD := New()

	bsA.b = append(bsA.b, 0b00000101, 0b10100000)
	bsB.b = append(bsB.b, 0b11000000, 0b00000011)
	bsC.b = append(bsC.b, 0b00000000, 0b00001000, 0b00000011)
	bsD.b = append(bsD.b, 0b00010000)

	bsA.Or(bsB)
	assert.Equal([]byte{0b11000101, 0b10100011}, bsA.b)

	bsA.Or(bsC)
	assert.Equal([]byte{0b11000101, 0b10101011}, bsA.b)

	bsA.Or(bsD)
	assert.Equal([]byte{0b11010101, 0b10101011}, bsA.b)
}
