package testsync

import "github.com/oligoden/chassis/lateral/interactor/query"

type TestSync struct {
	queries []query.Query
	index   int
}

func New(index int) *TestSync {
	return &TestSync{}
}

func (c *TestSync) Sync(qi []query.Query) []query.Query {
	return query.Sync(qi, c.queries, c.index)
}
