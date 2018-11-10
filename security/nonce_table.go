package security

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

type Nonce []byte

type NonceTable struct {
	nonceList map[byte]Nonce
	lock      *sync.Mutex
	timers    map[byte]*time.Timer
}

func NewNonceTable() *NonceTable {
	return &NonceTable{
		nonceList: map[byte]Nonce{},
		lock:      &sync.Mutex{},
		timers:    map[byte]*time.Timer{},
	}
}

func (t *NonceTable) Set(key byte, nonce Nonce, timeout time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.set(key, nonce, timeout)
}

func (t *NonceTable) Get(key byte) (Nonce, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	return t.get(key)
}

// @todo determine whether this is needed
// func (t *NonceTable) Peek(key byte) (Nonce, error) {
// 	t.lock.Lock()
// 	defer t.lock.Unlock()
//
// 	return t.peek(key)
// }

func (t *NonceTable) Delete(key byte) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.delete(key)
}

func (t *NonceTable) Generate(timeout time.Duration) (Nonce, error) {
	t.lock.Lock()
	defer t.lock.Unlock()

	for i := 0; ; i++ {
		nonce := GenerateNonce()

		// !ok indicates the item was not in the map, so we're good to save it and break
		if _, ok := t.nonceList[nonce[0]]; !ok {
			t.set(nonce[0], nonce, timeout)
			return nonce, nil
		}

		if i > 5 {
			return nil, errors.New("Unable to resolve nonce collision (are nonces being deleted?)")
		}
	}
}

func (t *NonceTable) set(key byte, nonce Nonce, timeout time.Duration) {
	// ensure the key does not exist in the map

	t.nonceList[key] = nonce

	if timer, ok := t.timers[key]; ok {
		timer.Reset(timeout)
	} else {
		t.timers[key] = time.AfterFunc(timeout, func() {
			t.Delete(key)
		})
	}
}

func (t *NonceTable) get(key byte) (Nonce, error) {
	nonce, err := t.peek(key)
	if err != nil {
		return nil, err
	}

	t.delete(key)

	return nonce, nil
}

func (t *NonceTable) peek(key byte) (Nonce, error) {
	if nonce, ok := t.nonceList[key]; ok {
		return nonce, nil
	}

	return nil, errors.New("key not found")
}

func (t *NonceTable) delete(key byte) {
	delete(t.nonceList, key)

	if timer, ok := t.timers[key]; ok {
		timer.Stop()
		delete(t.timers, key)
	}
}

func GenerateNonce() []byte {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		// @todo
		panic(err)
	}

	return buf
}
