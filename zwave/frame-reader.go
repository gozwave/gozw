package zwave

import "bufio"

func readFrames(reader *bufio.Reader, frames chan<- *ZFrame) {
	for {
		header, err := reader.ReadByte()
		if err != nil {
			// @todo handle panic
			panic(err)
		}

		if header != 0x01 {
			frame := UnmarshalFrame([]byte{header})
			frames <- frame
			continue
		}

		length, err := reader.ReadByte()
		if err != nil {
			// @todo handle panic
			panic(err)
		}

		buf := make([]byte, length+2)
		buf[0] = header
		buf[1] = length

		for i := 0; i < int(length)-1; i++ {
			data, err := reader.ReadByte()
			if err != nil {
				// @todo handle panic
				panic(err)
			}

			buf[i+2] = data
		}

		checksum, err := reader.ReadByte()
		if err != nil {
			// @todo handle panic
			panic(err)
		}

		buf[len(buf)-1] = checksum

		frame := UnmarshalFrame(buf)
		frames <- frame
	}
}
