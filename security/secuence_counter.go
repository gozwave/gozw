package security

import "sync"

const (
	SecuritySequenceCounterMin byte = 1
	SecuritySequenceCounterMax      = 15
)

type SequenceCounter struct {
	// maps a node id to a sequence counter (unique per node)
	counters map[byte]byte
	lock     *sync.Mutex
}

func NewSequenceCounter() *SequenceCounter {
	return &SequenceCounter{
		counters: map[byte]byte{},
		lock:     &sync.Mutex{},
	}
}

func (s *SequenceCounter) Get(nodeID byte) (counter byte) {
	var ok bool

	s.lock.Lock()
	defer s.lock.Unlock()

	if counter, ok = s.counters[nodeID]; !ok {
		s.counters[nodeID] = SecuritySequenceCounterMin
		return SecuritySequenceCounterMin
	}

	if counter+1 > SecuritySequenceCounterMax {
		counter = SecuritySequenceCounterMin
	} else {
		counter++
	}

	s.counters[nodeID] = counter

	return
}
