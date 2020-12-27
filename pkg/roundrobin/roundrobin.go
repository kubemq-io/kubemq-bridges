package roundrobin

import "go.uber.org/atomic"

type RoundRobin struct {
	index  *atomic.Int32
	length int
}

func NewRoundRobin(length int) *RoundRobin {
	rr := &RoundRobin{
		index:  atomic.NewInt32(0),
		length: length,
	}
	return rr
}

func (rr *RoundRobin) Next() int {
	next := int(rr.index.Load())
	if next+1 >= rr.length {
		rr.index.Store(0)
	} else {
		rr.index.Inc()
	}
	return next
}
