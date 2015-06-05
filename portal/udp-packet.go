package portal

import (
	"bytes"
	"encoding/binary"
	"net"
)

const udpHeaderLen = 8
const udpProtoNum = 0x11

type UDPv6Packet struct {
	SrcPort  uint16
	DstPort  uint16
	Checksum uint16
	Payload  []byte
}

func NewUDPv6Packet(srcPort, dstPort uint16, srcAddr, dstAddr net.IP, payload []byte) *UDPv6Packet {
	packet := UDPv6Packet{
		SrcPort: srcPort,
		DstPort: dstPort,
		Payload: payload,
	}

	packet.SetChecksum(srcAddr, dstAddr)

	return &packet
}

func (udp *UDPv6Packet) Marshal() []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.BigEndian, uint16(udp.SrcPort))
	binary.Write(buf, binary.BigEndian, uint16(udp.DstPort))
	binary.Write(buf, binary.BigEndian, uint16(udpHeaderLen+len(udp.Payload)))
	binary.Write(buf, binary.BigEndian, uint16(udp.Checksum))

	return append(buf.Bytes(), udp.Payload...)
}

// See https://github.com/grahamking/latency/blob/master/tcp.go#L120
func (udp *UDPv6Packet) SetChecksum(src, dst net.IP) {

	var csum uint32

	for i := 0; i < 16; i += 2 {
		csum += uint32(src[i]) << 8
		csum += uint32(src[i+1])
		csum += uint32(dst[i]) << 8
		csum += uint32(dst[i+1])
	}

	length := uint32(len(udp.Payload) + udpHeaderLen)

	csum += uint32(udpProtoNum) // next header is always UDP
	csum += length & 0xffff
	csum += length >> 16

	dataLen := len(udp.Payload) - 1

	for i := 0; i < dataLen; i += 2 {
		csum += uint32(udp.Payload[i]) << 8
		csum += uint32(udp.Payload[i+1])
	}

	if len(udp.Payload)%2 == 1 {
		csum += uint32(udp.Payload[dataLen]) << 8
	}

	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}

	// Bitwise complement
	udp.Checksum = ^uint16(csum + (csum >> 16))
}
