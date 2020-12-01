package gosql

import "strings"

type Join struct {
	joins []string
}

func NewJoin(js ...string) *Join {
	return &Join{
		joins: js,
	}
}

func (j *Join) Add(js ...string) {
	j.joins = append(j.joins, js...)
}

func (j *Join) Compile() string {
	return " " + strings.Join(j.joins, " ")
}
