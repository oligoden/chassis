package lateral

type Engin struct {
	enlistings chan enlistOp
	pushes     chan pushOp
	pulls      chan pullOp
	execs      chan execOp
}

type enlistOp struct {
	query Query
	resp  chan Query
}

type pushOp struct {
	query Query
}

type pullOp struct {
	curSeq uint
	resp   chan []Query
}

type execOp struct {
	resp chan Query
}

func NewEngin() Engin {
	return Engin{
		enlistings: make(chan enlistOp),
		pushes:     make(chan pushOp),
		pulls:      make(chan pullOp),
		execs:      make(chan execOp),
	}
}

func (e Engin) Interactor(i Interactor) {

}

func (e Engin) EnlistOp(query Query) Query {
	enlist := enlistOp{
		query: query,
		resp:  make(chan Query),
	}
	e.enlistings <- enlist
	return <-enlist.resp
}

func (e Engin) PushOp(query Query) {
	push := pushOp{
		query: query,
	}
	e.pushes <- push
}

func (e Engin) PullOp(curSeq uint) []Query {
	pull := pullOp{
		curSeq: curSeq,
		resp:   make(chan []Query),
	}
	e.pulls <- pull
	return <-pull.resp
}

func (e Engin) ExecOp() Query {
	exec := execOp{
		resp: make(chan Query),
	}
	e.execs <- exec
	return <-exec.resp
}

func (e Engin) Start() {
	var base, current uint
	qstore := []Query{}
	qqueue := map[uint]string{}

	for {
		select {
		case push := <-e.pushes:
			qqueue[push.query.Seq] = push.query.Body
		// if push.query.Seq - current > 1 {
		// 	ita.Pull(current)
		// }
		case pull := <-e.pulls:
			if pull.curSeq >= base+uint(len(qstore)) {
				pull.resp <- []Query{}
				break
			}
			from := base + pull.curSeq
			to := from + 10
			if to >= uint(len(qstore)) {
				to = uint(len(qstore))
			}
			pull.resp <- qstore[from:to]
		case exec := <-e.execs:
			q, exist := qqueue[current+1]
			if !exist {
				exec.resp <- Query{}
				break
			}
			query := Query{
				Seq:  current + 1,
				Body: q,
			}
			qstore = append(qstore, query)
			delete(qqueue, current)
			current++
			exec.resp <- query
		}
	}
}

type Interactor interface {
	Pull() []Query
}

type Query struct {
	Seq  uint
	Body string
}
