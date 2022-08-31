package interactor

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/oligoden/chassis/lateral/bitset"
	"github.com/oligoden/chassis/lateral/interactor/query"
)

type Interactor struct {
	index     int
	current   int
	Address   string
	Instances []Syncer
	// queue     []query.Query
	confirmed []query.Query
	Adds      chan addOp
	Updates   chan updateOp
}

// type query struct {
// 	UC       string
// 	tsp      time.Time
// 	instance string
// 	Query    interface{}
// 	Confirms *bitset.Seq
// }

// type queryList []query

// func (qs queryList) Len() int {
// 	return len(qs)
// }
// func (qs queryList) Swap(i, j int) {
// 	qs[i], qs[j] = qs[j], qs[i]
// }
// func (qs queryList) Less(i, j int) bool {
// 	return qs[i].tsp.Before(qs[j].tsp)
// }

// type timecode struct {
// 	UC       string
// 	tsp      time.Time
// 	instance int
// }

// An implementation of Syncer provides the Sync method for synchronising
// queues across service instances.
type Syncer interface {
	Sync([]query.Query) []query.Query
}

type addOp struct {
	query query.Query
}

type updateOp struct {
	queueIn  []query.Query
	queueOut chan []query.Query
}

func New(index int, address string) *Interactor {
	ita := &Interactor{
		index:     index,
		Address:   address,
		confirmed: []query.Query{},
		// syncer:    cm,
	}

	ita.Instances = make([]Syncer, index+1)
	ita.Instances[index] = ita

	ita.Adds = make(chan addOp)
	ita.Updates = make(chan updateOp)

	go func() {
		queue := []query.Query{}

		for {
			select {
			case add := <-ita.Adds:
				queue = append(queue, add.query)
			case upd := <-ita.Updates:
				queue = query.Sync(upd.queueIn, queue, ita.index)
				upd.queueOut <- queue
			case <-time.After(1 * time.Millisecond):
				if len(queue) > 0 {
					if queue[0].Confirms.Zero() {
						// all remotes aggreed, move to confirmed list
						ita.confirmed = append(ita.confirmed, queue[0])
						queue = queue[1:]
						ita.current++
						continue
					}

					for i := 0; i < queue[0].Confirms.Len(); i++ {
						if queue[0].Confirms.Get(i) {
							queue = query.Sync(ita.Instances[i].Sync(queue), queue, ita.index)
							break
						}
					}
				}
			}
		}
	}()

	return ita
}

func (ita *Interactor) Queue(qry interface{}) {
	q := query.Query{
		TS:       time.Now(),
		Instance: ita.Address,
		Query:    qry,
		Confirms: bitset.New(),
	}

	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%v", q)))
	q.UC = fmt.Sprintf("%x", h.Sum(nil))

	for i := range ita.Instances {
		q.Confirms.Set(i, true)
	}
	q.Confirms.Set(ita.index, false)

	add := addOp{
		query: q,
	}
	ita.Adds <- add
}

func (ita *Interactor) Sync(qi []query.Query) []query.Query {
	update := updateOp{
		queueIn:  qi,
		queueOut: make(chan []query.Query),
	}
	ita.Updates <- update
	return <-update.queueOut
}

func (ita *Interactor) Confirmed() []query.Query {
	return ita.confirmed
}
