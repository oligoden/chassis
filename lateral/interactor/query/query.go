package query

import (
	"time"

	"github.com/oligoden/chassis/lateral/bitset"
)

type Query struct {
	UC       string
	TS       time.Time
	Instance string
	Query    interface{}
	Confirms *bitset.Seq
}

func Sync(qi, qc []Query, ii int) (qo []Query) {
	exist := false

	for _, qIncoming := range qi {
		for _, qCurrent := range qc {
			if qIncoming.UC == qCurrent.UC {
				qCurrent.Confirms.Or(qIncoming.Confirms)
				exist = true
			}
		}

		if !exist {
			qc = append(qc, Query{
				UC:       qIncoming.UC,
				TS:       qIncoming.TS,
				Instance: qIncoming.Instance,
				Query:    qIncoming.Query,
				Confirms: qIncoming.Confirms,
			})

			qc[len(qc)-1].Confirms.Set(ii, false)
		}
	}

	qo = make([]Query, len(qc))
	copy(qo, qc)
	return qo
}
