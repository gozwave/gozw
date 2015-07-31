package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSetNonce(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	table.Set(0x00, []byte{0x00, 0x01, 0x02, 0x03}, time.Second)

	nonce, err := table.Get(0x00)
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x00, 0x01, 0x02, 0x03}, nonce)

	nonce, err = table.Get(0x01)
	assert.Error(t, err)
	assert.Nil(t, nonce)
}

func TestGetDeletesItem(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	table.Set(0x00, []byte{0x00, 0x01, 0x02, 0x03}, time.Second)

	nonce, err := table.Get(0x00)
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x00, 0x01, 0x02, 0x03}, nonce)

	nonce, err = table.Get(0x00)
	assert.Error(t, err)
	assert.Nil(t, nonce)
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	table.Set(0x00, []byte{0x00, 0x01, 0x02, 0x03}, time.Second)
	table.Delete(0x00)
	nonce, err := table.Get(0x00)

	assert.Error(t, err)
	assert.Nil(t, nonce)
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	nonce, err := table.Generate(time.Second)
	assert.NoError(t, err)
	assert.Len(t, nonce, 8)
}

func TestNoncesTimeOut(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	table.Set(0x00, []byte{0x00}, time.Microsecond*10)
	time.Sleep(time.Millisecond * 50)

	nonce, err := table.Get(0x00)
	assert.Error(t, err)
	assert.Nil(t, nonce)
}

func TestNonceSetResetsTimeout(t *testing.T) {
	t.Parallel()

	table := NewNonceTable()

	table.Set(0x00, []byte{0x00}, time.Millisecond*100)
	time.Sleep(time.Millisecond * 50)

	table.Set(0x00, []byte{0x01}, time.Millisecond*250)
	time.Sleep(time.Millisecond * 100)

	nonce, err := table.Get(0x00)
	assert.NoError(t, err)
	assert.EqualValues(t, []byte{0x01}, nonce)

	table.Set(0x00, []byte{0x00}, time.Millisecond*100)
	time.Sleep(time.Millisecond * 50)

	table.Set(0x00, []byte{0x01}, time.Millisecond*250)
	time.Sleep(time.Millisecond * 300)

	nonce, err = table.Get(0x00)
	assert.Error(t, err)
	assert.Nil(t, nonce)
}
