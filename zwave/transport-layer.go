package zwave

import (
	"bufio"
	"io"

	"github.com/tarm/serial"
)

type TransportLayer interface {
	Read() <-chan byte
	Write(buf []byte) (int, error)
}

type SerialTransportLayer struct {
	port *serial.Port
}

func NewSerialTransportLayer(device string, baud int) (*SerialTransportLayer, error) {
	var err error

	port, err := serial.OpenPort(&serial.Config{
		Name: device,
		Baud: baud,
	})

	if err != nil {
		return nil, err
	}

	transport := &SerialTransportLayer{
		port: port,
	}

	return transport, nil
}

func (t *SerialTransportLayer) Read() <-chan byte {
	readQueue := make(chan byte)

	go t.readAsync(readQueue)

	return readQueue
}

func (t *SerialTransportLayer) Write(buf []byte) (int, error) {
	return t.port.Write(buf)
}

func (t *SerialTransportLayer) readAsync(readQueue chan<- byte) {
	reader := bufio.NewReader(t.port)

	for {
		byt, err := reader.ReadByte()

		if err == io.EOF {
			close(readQueue)
			break
		} else if err != nil {
			panic(err)
		}

		readQueue <- byt
	}
}
