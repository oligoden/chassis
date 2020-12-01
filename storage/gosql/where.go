package gosql

import "fmt"

type Where struct {
	group     []Where
	condition string
	values    []interface{}
	and       bool
}

func NewWhere(q string, vs ...interface{}) *Where {
	return &Where{
		condition: q,
		values:    vs,
	}
}

func NewWhereGroup(wg *Where) *Where {
	return &Where{
		condition: fmt.Sprintf("(%s)", wg.condition),
		values:    wg.values,
	}
}

func (w *Where) And(q string, vs ...interface{}) *Where {
	w.condition = fmt.Sprintf("%s AND %s", w.condition, q)
	w.values = append(w.values, vs...)
	return w
}

func (w *Where) Or(q string, vs ...interface{}) *Where {
	w.condition = fmt.Sprintf("%s OR %s", w.condition, q)
	w.values = append(w.values, vs...)
	return w
}

func (w *Where) AndGroup(wg *Where) *Where {
	w.condition = fmt.Sprintf("%s AND (%s)", w.condition, wg.condition)
	w.values = append(w.values, wg.values...)
	return w
}

func (w *Where) OrGroup(wg *Where) *Where {
	w.condition = fmt.Sprintf("%s OR (%s)", w.condition, wg.condition)
	w.values = append(w.values, wg.values...)
	return w
}

func (w *Where) Compile() (string, []interface{}) {
	return " WHERE " + w.condition, w.values
}
