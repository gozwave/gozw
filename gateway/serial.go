package gateway

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/zwave"
	"github.com/tarm/serial"
)

type ZCallback func(response *zwave.ZFrame)

type Request struct {
	frame    *zwave.ZFrame
	callback ZCallback
}

type SerialPort struct {
	port            *serial.Port
	incomingPrivate chan *zwave.ZFrame
	requestQueue    chan Request
	Incoming        chan *zwave.ZFrame
}

func NewSerialPort(config *common.GozwConfig) (*SerialPort, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: config.Device,
		Baud: config.Baud,
	})

	if err != nil {
		return nil, err
	}

	requestQueue := make(chan Request, 1)
	incomingPrivate := make(chan *zwave.ZFrame, 1)
	incomingPublic := make(chan *zwave.ZFrame, 1)

	serialPort := SerialPort{
		port:            port,
		incomingPrivate: incomingPrivate,
		requestQueue:    requestQueue,
		Incoming:        incomingPublic,
	}

	return &serialPort, nil
}

func (s *SerialPort) Initialize() {
	err := s.sendNak()
	if err != nil {
		panic(err)
	}

	go readFrames(bufio.NewReader(s.port), s.incomingPrivate)

	for {
		select {
		case frame := <-s.incomingPrivate:
			go func(frame *zwave.ZFrame) {
				s.Incoming <- frame
			}(frame)
		case <-time.After(time.Second * 2):
			return
		}
	}
}

func (s *SerialPort) Run() {
	for {
		select {
		case incoming := <-s.incomingPrivate:
			s.Incoming <- incoming
		case request := <-s.requestQueue:
			s.transmitFrame(request.frame)
			s.sendAck()
			fmt.Println(<-s.incomingPrivate)
		}
	}
}

func (s *SerialPort) SendFrameSync(frame *zwave.ZFrame) *zwave.ZFrame {
	await := make(chan *zwave.ZFrame, 1)
	callback := func(response *zwave.ZFrame) {
		await <- response
	}

	go s.SendFrame(frame, callback)

	return <-await
}

func (s *SerialPort) SendFrame(frame *zwave.ZFrame, callback ZCallback) {
	go func(frame *zwave.ZFrame, callback ZCallback) {
		s.requestQueue <- Request{
			frame:    frame,
			callback: callback,
		}
	}(frame, callback)
}

func (s *SerialPort) Close() error {
	return s.port.Close()
}

func (s *SerialPort) sendAck() error {
	_, err := s.port.Write(zwave.NewAckFrame().Marshal())
	return err
}

func (s *SerialPort) sendNak() error {
	_, err := s.port.Write(zwave.NewNakFrame().Marshal())
	return err
}

func (s *SerialPort) transmitFrame(frame *zwave.ZFrame) *zwave.ZFrame {

	// Write the frame to the serial port
	numBytes, err := s.port.Write(frame.Marshal())
	if err != nil {
		panic(err)
	}

	fmt.Printf("---> %d bytes written:\n", numBytes)
	fmt.Println(hex.Dump(frame.Marshal()))

	// Block for the next incoming frame
	receipt := <-s.incomingPrivate

	// If the
	if receipt.IsAck() {
		receipt = <-s.incomingPrivate
		return receipt
	}

	fmt.Println("BAD FRAME!!!", receipt)
	panic("hi")
}
