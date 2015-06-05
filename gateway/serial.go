package gateway

import (
	"bufio"
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

// NewSerialPort Open a(n actual) serial port and create some supporting channels
func NewSerialPort(config *common.GozwConfig) (*SerialPort, error) {
	// Open the serial port with the given device and baud rate
	// Note: could probably consider inlining the baud rate, since it should
	// always be 115200
	port, err := serial.OpenPort(&serial.Config{
		Name: config.Device,
		Baud: config.Baud,
	})

	if err != nil {
		return nil, err
	}

	// Channel for Z-Wave commands we need to queue up
	requestQueue := make(chan Request, 1)

	// Channel for frames we receive from the Z-Wave controller, but don't necessarily
	// want to make public (yet), since this can include ACKs, NAKs, and CANs
	incomingPrivate := make(chan *zwave.ZFrame, 1)

	// Channel for frames we want to release into the wild
	incomingPublic := make(chan *zwave.ZFrame, 1)

	serialPort := SerialPort{
		port:            port,
		incomingPrivate: incomingPrivate,
		requestQueue:    requestQueue,
		Incoming:        incomingPublic,
	}

	return &serialPort, nil
}

// Initialize We need to do some initial setup on the device before we are able
// to enter our normal handler loop
func (s *SerialPort) Initialize() {

	// According to 6.1 in the Serial API guide, we're supposed to start up by
	// sending a NAK, then doing a hard or soft reset. Soft reset isn't implemented
	// yet, and I don't know if hard reset is possible with a USB controller
	err := s.sendNak()
	if err != nil {
		panic(err)
	}

	// Read frames from the serial port in a goroutine, and make them available on the
	// incomingPrivate channel
	go readFrames(bufio.NewReader(s.port), s.incomingPrivate)

	// This block will block to receive incoming frames and continue to do so until
	// 2 seconds after it has received the last frame. We do this because if we
	// previously crashed, a quick startup could bring us up while the controller is
	// still retransmitting frames we haven't ACKed, and we might not know what to do
	// with them
	for {
		select {
		case frame := <-s.incomingPrivate:
			// this runs in a goroutine in case nothing is listening to s.Incoming yet
			// the goroutine basically just blocks until something listens.
			go func(frame *zwave.ZFrame) {
				s.Incoming <- frame
			}(frame)
		case <-time.After(time.Second * 2):
			// after 2 seconds of not receiving any frames, return
			return
		}
	}
}

// Run Handles unsolicited incoming frames and transmits outgoing frames queued
// using SendFrame
func (s *SerialPort) Run() {
	for {
		select {
		case incoming := <-s.incomingPrivate:
			// @todo verify checksum and send appropriate response frame
			s.Incoming <- incoming
		case request := <-s.requestQueue:
			s.transmitFrame(request.frame)
			s.sendAck()
			fmt.Println(<-s.incomingPrivate)
		}
	}
}

// SendFrameSync wraps SendFrame with some magic that blocks until the result
// arrives
func (s *SerialPort) SendFrameSync(frame *zwave.ZFrame) *zwave.ZFrame {
	// Make a channel we can block on
	await := make(chan *zwave.ZFrame, 1)

	// All our callback needs to do is publish the response frame back to the channel
	callback := func(response *zwave.ZFrame) {
		await <- response
	}

	// Send the frame in a goroutine, since we don't want to block on this
	go s.SendFrame(frame, callback)

	// Block until the channel emits a value for us, and then return that value
	return <-await
}

// SendFrame queues a frame to be sent to the controller
func (s *SerialPort) SendFrame(frame *zwave.ZFrame, callback ZCallback) {
	go func(frame *zwave.ZFrame, callback ZCallback) {
		s.requestQueue <- Request{
			frame:    frame,
			callback: callback,
		}
	}(frame, callback)
}

// Close the serial port
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
	_, err := s.port.Write(frame.Marshal())
	if err != nil {
		panic(err)
	}

	// fmt.Printf("---> %d bytes written:\n", numBytes)
	// fmt.Println(hex.Dump(frame.Marshal()))

	// Block for the next incoming frame
	// @todo possible race condition here, since we block on s.incomingPrivate in Run()
	// 	     should refactor SerialPort into a state machine that stores the request
	//       in flight. We can handle the ACK/NAK/CAN response there, as well as any
	//       response frames
	receipt := <-s.incomingPrivate

	// If the next frame is an ACK, we can wait for the next frame, which should
	// be a response frame
	if receipt.IsAck() {
		receipt = <-s.incomingPrivate
		return receipt
	}

	// @todo Handle non-ACK frames here
	// @todo Not all frames receive a response frame, so it isn't appropriate to
	// always grab the next incoming frame.

	// @todo Need to reinitialize our connection and do a reset of the controller
	fmt.Println("BAD FRAME!!!", receipt)
	panic("hi")
}
