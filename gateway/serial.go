package gateway

import (
	"bufio"
	"fmt"
	"time"

	"github.com/bjyoungblood/gozw/common"
	"github.com/bjyoungblood/gozw/zwave"
	"github.com/bjyoungblood/gozw/zwave/layers"
	"github.com/google/gopacket"
	"github.com/tarm/serial"
)

// AckCallback is a function callback to be executed when a frame is transmitted.
// status will be one of zwave.FrameHeader*
type AckCallback func(responseType uint8, response *zwave.Frame)

// Request represents a ZFrame queued for transmission to the controller
type Request struct {
	frame    *zwave.Frame
	callback AckCallback
	attempts int
}

// SerialPort is a container/wrapper for the actual serial port, with some
// extra protection to ensure proper connection state with the controller
type SerialPort struct {
	port *serial.Port

	// Channel for parsed frames (packets)
	incomingPackets chan gopacket.Packet

	// Channel for Z-Wave commands we need to queue up
	requestQueue chan Request

	// Storage for the currently-running request
	requestInFlight Request

	// Channel for frames we want to release into the wild
	Incoming chan *zwave.Frame
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

	serialPort := SerialPort{
		port: port,

		incomingPackets: make(chan gopacket.Packet, 1),
		requestQueue:    make(chan Request, 1),
		Incoming:        make(chan *zwave.Frame, 1),
	}

	return &serialPort, nil
}

func (s *SerialPort) ReadPacketData() ([]byte, gopacket.CaptureInfo, error) {
	buf := make([]byte, 128)
	readLen, err := s.port.Read(buf)

	ci := gopacket.CaptureInfo{
		Timestamp:     time.Now(),
		CaptureLength: readLen,
		Length:        readLen,
	}

	return buf, ci, err
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
	// incomingPackets channel
	go readFrames(s.port, s.incomingPackets)

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
		// case <-s.incomingPackets:
		// 	// this runs in a goroutine in case nothing is listening to s.Incoming yet
		// 	// the goroutine basically just blocks until something listens.
		// 	// go func(packet gopacket.Packet) {
		// 	// 	s.Incoming <- packet
		// 	// }(packet)
		case <-time.After(time.Second * 2):
			// after 2 seconds of not receiving any frames, return
			continue
		}
	}
}

// Run Handles unsolicited incoming frames and transmits outgoing frames queued
// using SendFrame
func (s *SerialPort) Run() {
	for {
		select {
		// case incoming := <-s.incomingPackets:
		// 	err := incoming.VerifyChecksum()
		// 	if err != nil {
		// 		s.sendNak()
		// 		continue
		// 	} else if incoming.IsData() {
		// 		// If everything else has been processed, then release it into the wild
		// 		s.sendAck()
		// 		s.Incoming <- incoming
		// 	} else {
		// 		fmt.Println("Unexpected frame: ", incoming)
		// 	}
		//
		// case request := <-s.requestQueue:
		// 	s.requestInFlight = request
		// 	_, err := s.port.Write(request.frame.Marshal())
		// 	if err != nil {
		// 		panic(err)
		// 	}
		//
		// 	confirmation := <-s.incomingPackets
		//
		// 	if confirmation.IsNak() || confirmation.IsCan() {
		// 		s.requestInFlight.callback(confirmation.Header, nil)
		// 	} else if confirmation.IsAck() {
		//
		// 		response := <-s.incomingPackets
		//
		// 		if response.IsData() {
		// 			s.sendAck()
		// 		}
		//
		// 		go s.requestInFlight.callback(confirmation.Header, response)
		// 	}

		// time.Sleep(10 * time.Millisecond)
		}
	}
}

// SendFrameSync wraps SendFrame with some magic that blocks until the result
// arrives
func (s *SerialPort) SendFrameSync(frame *zwave.Frame) *zwave.Frame {
	// Make a channel we can block on
	await := make(chan *zwave.Frame, 1)

	// All our callback needs to do is publish the response frame back to the channel
	callback := func(response uint8, responseFrame *zwave.Frame) {
		await <- responseFrame
	}

	// Send the frame in a goroutine, since we don't want to block on this
	go s.SendFrame(frame, callback)

	// Block until the channel emits a value for us, and then return that value
	return <-await
}

// SendFrame queues a frame to be sent to the controller
func (s *SerialPort) SendFrame(frame *zwave.Frame, callback AckCallback) {
	go func(frame *zwave.Frame, callback AckCallback) {
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

func (s *SerialPort) Write(buf []byte) (int, error) {
	written, err := s.port.Write(buf)
	return written, err
}

func (s *SerialPort) sendAck() error {
	_, err := s.port.Write(zwave.NewAckFrame().Marshal())
	return err
}

func (s *SerialPort) sendNak() error {
	_, err := s.port.Write(zwave.NewNakFrame().Marshal())
	return err
}

// @todo handle EOF, other errors instead of panic
func readFrames(port *serial.Port, incomingPackets chan<- gopacket.Packet) {
	reader := bufio.NewReader(port)

	for {
		// Read the SOF byte
		sof, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}

		// Handle ACK, CAN, and NAK frames first
		if sof == layers.FrameSOFAck || sof == layers.FrameSOFCan || sof == layers.FrameSOFNak {
			packet := gopacket.NewPacket([]byte{sof}, layers.LayerTypeFrame, gopacket.DecodeOptions{})
			incomingPackets <- packet
			continue
		}

		// If we're seeing something other than a data SOF here, we need to ignore it
		// to flush garbage out of the read buffer, per specification
		if sof != layers.FrameSOFData {
			continue
		}

		// Read the length from the frame
		length, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}

		buf := make([]byte, length+2)
		buf[0] = sof
		buf[1] = length

		// read the frame payload
		for i := 0; i < int(length)-1; i++ {
			data, err := reader.ReadByte()
			if err != nil {
				// @todo handle panic
				panic(err)
			}

			buf[i+2] = data
		}

		// read the checksum
		checksum, err := reader.ReadByte()
		if err != nil {
			// @todo handle panic
			panic(err)
		}

		buf[len(buf)-1] = checksum

		packet := gopacket.NewPacket(buf, layers.LayerTypeFrame, gopacket.DecodeOptions{})
		fmt.Println(packet.Dump())
		incomingPackets <- packet
	}
}
