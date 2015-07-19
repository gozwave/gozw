package mocks

import "github.com/stretchr/testify/mock"

type TransportLayer struct {
	mock.Mock
}

func (m *TransportLayer) Read() <-chan byte {
	ret := m.Called()

	var r0 <-chan byte
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(<-chan byte)
	}

	return r0
}

func (m *TransportLayer) Write(buf []byte) (int, error) {
	ret := m.Called(buf)

	r0 := ret.Get(0).(int)
	r1 := ret.Error(1)

	return r0, r1
}
