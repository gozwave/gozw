package frame

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockedTransportLayer struct {
	mock.Mock

	input chan byte
}

func (t *MockedTransportLayer) Write(bytes []byte) (int, error) {
	args := t.Called(bytes)
	written := args.Get(0).(int)
	err := args.Get(1).(error)

	return written, err
}

func (t *MockedTransportLayer) Read() <-chan byte {
	args := t.Called()
	return args.Get(0).(chan byte)
}

func TestCanary(t *testing.T) {
	// frameLayer := NewFrameLayer(&MockedTransportLayer{})
	// output := frameLayer.GetOutputChannel()
}
