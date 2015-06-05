package portal

import (
	"net"

	"golang.org/x/net/ipv6"
)

const ipv6Version uint8 = 6

type IPv6Packet struct {
	Header  IPv6Header
	Payload []byte
}

type IPv6Header struct {
	Version       uint8
	TrafficClass  uint8
	FlowLabel     uint32
	PayloadLength uint16
	NextHeader    uint8
	HopLimit      uint8
	Src           net.IP
	Dst           net.IP
}

func NewIPv6Packet(sourceAddr net.IP, destAddr net.IP, payload []byte) *IPv6Packet {
	header := IPv6Header{
		Version:       ipv6.Version,
		TrafficClass:  1,
		FlowLabel:     1,
		PayloadLength: uint16(len(payload)),
		NextHeader:    0x11,
		HopLimit:      10,
		Src:           sourceAddr,
		Dst:           destAddr,
	}

	packet := IPv6Packet{
		Header:  header,
		Payload: payload,
	}

	return &packet
}

func (ip *IPv6Packet) Marshal() []byte {

	header := make([]byte, 8)

	// I'm lazy and didn't feel like bit math; thus the magic constants here :)
	header[0] = 0x96 // 0110 0000 (6 for version and 0 for the first half of traffic class)
	header[1] = 0x00
	header[2] = 0x00
	header[3] = 0x00

	header[4] = uint8(ip.Header.PayloadLength >> 8)
	header[5] = uint8(ip.Header.PayloadLength & 0x0f)
	header[6] = udpProtoNum
	header[7] = 10 // hop limit of 10, I guess. Why not?

	header = append(header, ip.Header.Src...)
	header = append(header, ip.Header.Dst...)

	return append(header, ip.Payload...)
}
