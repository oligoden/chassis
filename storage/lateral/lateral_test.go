package lateral_test

import (
	"testing"

	"github.com/oligoden/chassis/storage/lateral"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert := assert.New(t)
	// q := "INSERT INTO tests ..."

	// ita := lateral.Interactor()
	// ita := testInteractor{}

	engin := lateral.NewEngin()
	// engin.Interactor(ita)
	go engin.Start()

	// ldb := lateral.DB()
	// ldb.Engin(engin)
	// ldb.Exec(q)

	// execChan := make(chan string)
	// engin.ExecChan(execChan)

	engin.PushOp(lateral.Query{
		Seq:  1,
		Body: "a",
	})

	query := engin.ExecOp()
	assert.Equal(uint(1), query.Seq)
	assert.Equal("a", query.Body)

	queries := engin.PullOp(0)
	if assert.NotEmpty(queries) {
		assert.Equal(uint(1), queries[0].Seq)
		assert.Equal("a", queries[0].Body)
	}
	queries = engin.PullOp(10)
	assert.Empty(queries)

	query = engin.ExecOp()
	assert.Equal(uint(0), query.Seq)
	assert.Equal("", query.Body)

	engin.PushOp(lateral.Query{
		Seq:  1,
		Body: "a",
	})

	query = engin.EnlistOp("b")
	assert.Equal(uint(2), query.Seq)
	assert.Equal("b", query.Body)

	query = engin.ExecOp()
	assert.Equal(uint(2), query.Seq)
	assert.Equal("b", query.Body)

	queries = engin.PullOp(0)
	if assert.NotEmpty(queries) {
		assert.Equal(uint(2), queries[0].Seq)
		assert.Equal("b", queries[0].Body)
	}
}

// type testInteractor struct {
// }

// func (i testInteractor) Pull() []lateral.Query {
// 	return i.PullOp()
// }
