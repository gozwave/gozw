package zwave

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

const (
	StateUninitialized = iota
	StateReady
	StateAwaitingAck
)

type ZCallback func(response *ZFrame)

type Request struct {
	frame    *ZFrame
	callback ZCallback
}

type SerialConfig struct {
	Device string
	Baud   int
}

type SerialPort struct {
	port            *serial.Port
	state           int
	incomingPrivate chan *ZFrame
	requestQueue    chan Request
	Incoming        chan *ZFrame
}

func NewSerialPort(config *SerialConfig) (*SerialPort, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: config.Device,
		Baud: config.Baud,
	})

	if err != nil {
		return nil, err
	}

	requestQueue := make(chan Request, 1)
	incomingPrivate := make(chan *ZFrame, 1)
	incomingPublic := make(chan *ZFrame, 1)

	serialPort := SerialPort{
		port:            port,
		state:           StateUninitialized,
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
			go func(frame *ZFrame) {
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

func (s *SerialPort) SendFrameSync(frame *ZFrame) *ZFrame {
	await := make(chan *ZFrame, 1)
	callback := func(response *ZFrame) {
		await <- response
	}

	go s.SendFrame(frame, callback)

	return <-await
}

func (s *SerialPort) SendFrame(frame *ZFrame, callback ZCallback) {
	go func(frame *ZFrame, callback ZCallback) {
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
	_, err := s.port.Write(NewAckFrame().Marshal())
	return err
}

func (s *SerialPort) transmitFrame(frame *ZFrame) *ZFrame {

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

func (s *SerialPort) sendNak() error {
	_, err := s.port.Write(NewNakFrame().Marshal())
	return err
}
