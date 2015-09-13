package transport

import (
	"bufio"

	"github.com/tarm/serial"
)

type SerialPortTransport struct {
	port   *serial.Port
	reader *bufio.Reader
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
		port:   port,
		reader: bufio.NewReader(port),
	}

	return transport, nil
}

func (t *SerialPortTransport) Read(p []byte) (n int, err error) {
	return t.reader.Read(p)
}

func (t *SerialPortTransport) ReadByte() (byte, error) {
	return t.reader.ReadByte()
}

func (t *SerialPortTransport) Write(buf []byte) (int, error) {
	return t.port.Write(buf)
}
