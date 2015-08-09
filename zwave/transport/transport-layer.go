package transport

import (
	"bufio"
	"io"

	"github.com/tarm/serial"
)

type Transport interface {
	Read() <-chan byte
	Write(buf []byte) (int, error)
}

type SerialPortTransport struct {
	port *serial.Port
}

func NewSerialPortTransport(device string, baud int) (*SerialPortTransport, error) {
	var err error

	port, err := serial.OpenPort(&serial.Config{
		Name: device,
		Baud: baud,
	})

	if err != nil {
		return nil, err
	}

	transport := &SerialPortTransport{
		port: port,
	}

	return transport, nil
}

func (t *SerialPortTransport) Read() <-chan byte {
	readQueue := make(chan byte)

	go t.readAsync(readQueue)

	return readQueue
}

func (t *SerialPortTransport) Write(buf []byte) (int, error) {
	return t.port.Write(buf)
}

func (t *SerialPortTransport) readAsync(readQueue chan<- byte) {
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
