package bitset

type Controller interface {
	Get(int) bool
	Set(int, bool)
	Len(...int) int
}

type Seq struct {
	b []byte
	l int
}

func New() *Seq {
	bs := &Seq{}
	bs.b = []byte{}
	bs.l = 0
	return bs
}

func (bs *Seq) Get(index int) bool {
	i := index % 8
	y := (index - i) / 8
	return bs.b[y]>>i&0b1 == 1
}

func (bs *Seq) Set(index int, b bool) {
	i := index % 8
	y := (index - i) / 8

	for {
		if y+1 > len(bs.b) {
			bs.b = append(bs.b, 0)
		} else {
			break
		}
	}

	if b {
		bs.b[y] = bs.b[y] | 1<<i
	} else {
		bs.b[y] = bs.b[y] ^ 1<<i&bs.b[y]
	}

	if bs.l < index+1 {
		bs.l = index + 1
	}
}

func (bs *Seq) Len(ls ...int) int {
	if len(ls) > 0 {
		index := ls[0] - 1
		i := index % 8
		y := (index - i) / 8

		if ls[0] < bs.l {
			bs.b = bs.b[0 : y+1]
			for j := i + 1; j < 8; j++ {
				bs.b[y] = bs.b[y] ^ 1<<j&bs.b[y]
			}
		}

		if ls[0] > bs.l {
			for {
				if y+1 > len(bs.b) {
					bs.b = append(bs.b, 0)
				} else {
					break
				}
			}
		}

		bs.l = ls[0]
	}

	return bs.l
}

func (bs *Seq) Zero() bool {
	for _, b := range bs.b {
		if b > 0 {
			return false
		}
	}

	return true
}

func (bsA *Seq) Or(bsB *Seq) {
	for i := range bsA.b {
		if i < len(bsB.b) {
			bsA.b[i] = bsA.b[i] | bsB.b[i]
		}
	}
}
