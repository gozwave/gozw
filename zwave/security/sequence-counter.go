package security

import (
	"runtime"
	"sync"
)

const (
	SecuritySequenceCounterMin uint8 = 1
	SecuritySequenceCounterMax       = 15
)

type SequenceCounter struct {
	// maps a node id to a sequence counter (unique per node)
	counters map[uint8]uint8
	lock     *sync.Mutex
}

func NewSequenceCounter() *SequenceCounter {
	return &SequenceCounter{
		counters: map[uint8]uint8{},
		lock:     &sync.Mutex{},
	}
}

func (s *SequenceCounter) Get(nodeId uint8) (counter uint8) {
	var ok bool

	s.lock.Lock()
	defer s.lock.Unlock()
	defer runtime.Gosched()

	if counter, ok = s.counters[nodeId]; !ok {
		s.counters[nodeId] = SecuritySequenceCounterMin
		return SecuritySequenceCounterMin
	}

	if counter+1 > SecuritySequenceCounterMax {
		counter = SecuritySequenceCounterMin
	} else {
		counter += 1
	}

	s.counters[nodeId] = counter

	return
}
