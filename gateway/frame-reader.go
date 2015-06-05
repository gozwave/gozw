package gateway

import (
	"bufio"

	"github.com/bjyoungblood/gozw/zwave"
)

// @todo handle EOF, other errors instead of panic
func readFrames(reader *bufio.Reader, frames chan<- *zwave.ZFrame) {
	// Loop forever
	for {
		// Read the header byte
		header, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}

		// If it's not a data frame, we don't need to read anymore
		// @todo should check that it is indeed a valid frame type; if not, we're
		// supposed to just discard it, as per the spec
		if header != 0x01 {
			frame := zwave.UnmarshalFrame([]byte{header})
			frames <- frame
			continue
		}

		// Read the length from the frame
		length, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}

		buf := make([]byte, length+2)
		buf[0] = header
		buf[1] = length

		// read the command id and payload
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

		// Unmarshal the byte array into a ZFrame
		frame := zwave.UnmarshalFrame(buf)
		frames <- frame
	}
}
