package transport

import (
	"bufio"

	"github.com/pkg/errors"
	"github.com/tarm/serial"
)

// SerialPortTransport contains a transport
type SerialPortTransport struct {
	port   *serial.Port
	reader *bufio.Reader
}

// NewSerialPortTransport will return a new serial port transport.
func NewSerialPortTransport(device string, baud int) (*SerialPortTransport, error) {
	var err error

	port, err := serial.OpenPort(&serial.Config{
		Name: device,
		Baud: baud,
	})
	if err != nil {
		return nil, errors.Wrap(err, "open port")
	}

	transport := &SerialPortTransport{
		port:   port,
		reader: bufio.NewReader(port),
	}

	return transport, nil
}

// Close will close the transport.
func (t *SerialPortTransport) Close() {
	t.port.Close()
}

// Read implements the io.Reader interface.
func (t *SerialPortTransport) Read(p []byte) (int, error) {
	return t.reader.Read(p)
}

// ReadByte implements the io.ByteReader interface.
func (t *SerialPortTransport) ReadByte() (byte, error) {
	return t.reader.ReadByte()
}

// Write implements the io.Writer interface.
func (t *SerialPortTransport) Write(buf []byte) (int, error) {
	return t.port.Write(buf)
}
