package gateway

import (
	"bufio"
	"time"

	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/zwave"
	"github.com/tarm/serial"
)

// AckCallback is a function callback to be executed when a frame is transmitted.
// status will be one of zwave.FrameHeader*
type AckCallback func(status int)

// Request represents a ZFrame queued for transmission to the controller
type Request struct {
	frame    *zwave.ZFrame
	callback AckCallback
	attempts int
}

// SerialPort is a container/wrapper for the actual serial port, with some
// extra protection to ensure proper connection state with the controller
type SerialPort struct {
	port            *serial.Port
	incomingPrivate chan *zwave.ZFrame
	requestQueue    chan Request
	requestInFlight Request
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

	// s.SendFrame(zwave.NewRequestFrame(zwave.ReadyCommand()), func(status int) {
	// 	fmt.Println(status)
	// })

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
			err := incoming.VerifyChecksum()
			if err != nil {
				s.sendNak()
				continue
			} else if incoming.IsAck() || incoming.IsNak() {
				if &s.requestInFlight != nil && s.requestInFlight.callback != nil {
					s.requestInFlight.callback(int(incoming.Header))
					s.requestInFlight = Request{}
					continue
				}
			} else if incoming.IsData() {
				// If everything else has been processed, then release it into the wild
				s.sendAck()
				s.Incoming <- incoming
			}

		case request := <-s.requestQueue:
			s.requestInFlight = request
			_, err := s.port.Write(request.frame.Marshal())
			if err != nil {
				panic(err)
			}

			// time.Sleep(10 * time.Millisecond)
		}
	}
}

// SendFrameSync wraps SendFrame with some magic that blocks until the result
// arrives
func (s *SerialPort) SendFrameSync(frame *zwave.ZFrame) int {
	// Make a channel we can block on
	await := make(chan int, 1)

	// All our callback needs to do is publish the response frame back to the channel
	callback := func(response int) {
		await <- response
	}

	// Send the frame in a goroutine, since we don't want to block on this
	go s.SendFrame(frame, callback)

	// Block until the channel emits a value for us, and then return that value
	return <-await
}

// SendFrame queues a frame to be sent to the controller
func (s *SerialPort) SendFrame(frame *zwave.ZFrame, callback AckCallback) {
	go func(frame *zwave.ZFrame, callback AckCallback) {
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
