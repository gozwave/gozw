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
	byteChannel := make(chan byte)

	go t.readAsync(byteChannel)

	return byteChannel
}

func (t *TransportLayer) Write(byteChannel []byte) (int, error) {
	return t.port.Write(byteChannel)
}

func (t *TransportLayer) readAsync(byteChannel chan<- byte) {
	reader := bufio.NewReader(t.port)

	for {
		byt, err := reader.ReadByte()

		if err == io.EOF {
			close(byteChannel)
			break
		} else if err != nil {
			panic(err)
		}

		byteChannel <- byt
	}
}
