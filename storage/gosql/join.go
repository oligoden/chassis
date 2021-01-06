package gosql

type Join struct {
	join string
}

func NewJoin(j string) *Join {
	return &Join{
		join: j,
	}
}

func (j *Join) Compile(ops ...string) (string, []interface{}) {
	return j.join, []interface{}{}
}

func (w *Join) Order() int {
	return 0
}
