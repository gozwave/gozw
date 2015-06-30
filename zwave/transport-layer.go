package zwave

import (
	"bufio"
	"io"

	"github.com/tarm/serial"
)

type TransportLayer struct {
	port *serial.Port
}

func NewTransportLayer(device string, baud int) (*TransportLayer, error) {
	var err error

	port, err := serial.OpenPort(&serial.Config{
		Name: device,
		Baud: baud,
	})

	if err != nil {
		return nil, err
	}

	transport := &TransportLayer{
		port: port,
	}

	return transport, nil
}

func (t *TransportLayer) Read() <-chan byte {
	readQueue := make(chan byte)

	go t.readAsync(readQueue)

	return readQueue
}

func (t *TransportLayer) Write(buf []byte) (int, error) {
	return t.port.Write(buf)
}

func (t *TransportLayer) readAsync(readQueue chan<- byte) {
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
