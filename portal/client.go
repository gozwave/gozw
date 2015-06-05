package portal

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"

	"github.com/bjyoungblood/gozw/zwave"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	UNINITIALIZED = iota
	GATEWAY_CONFIGURATION_SENT
	GATEWAY_CONFIGURATION_ACK
	INITIALIZED
)

// Client represents a connection from a Z/IP gateway
type Client struct {
	state int
	conn  net.Conn
	ch    chan []byte
}

func NewClient(conn net.Conn) Client {
	return Client{
		conn:  conn,
		state: UNINITIALIZED,
	}

}

func (c *Client) Handle() {

	c.conn.Write(zwave.MakeConfigurationSetPacket())
	fmt.Println("<--- CONFIGURATION_SET:")
	fmt.Println(hex.Dump(zwave.MakeConfigurationSetPacket()))
	fmt.Println("----")

	reader := bufio.NewReader(c.conn)

	buf := make([]byte, 64)
	byteCount, err := reader.Read(buf)
	if err != nil {
		fmt.Println("Error reading from client")
		c.Close()
		return
	}

	if byteCount != 3 {
		fmt.Printf("Unexpected byte count (%d) %s\n", byteCount, string(buf[0:byteCount]))
	}

	if buf[0] == 0x61 && buf[1] == 0x02 && buf[2] == 0xff {
		fmt.Println("---> CONFIGURATION_SET OK")
	} else {
		fmt.Printf("???? %s\n", buf)
	}

	serializeBuf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	ipLayer := layers.IPv6{
		Version:      6,
		TrafficClass: 0,
		FlowLabel:    0,
		NextHeader:   layers.IPProtocolUDP,
		HopLimit:     10,
		SrcIP:        net.ParseIP("3000::1"),
		DstIP:        net.ParseIP("ff02::2"),
	}

	udpLayer := layers.UDP{
		SrcPort: 44123,
		DstPort: 4123,
	}

	udpLayer.SetNetworkLayerForChecksum(&ipLayer)

	err = gopacket.SerializeLayers(serializeBuf, opts,
		&ipLayer,
		&udpLayer,
		gopacket.Payload(zwave.MakeNodeInfoCachedGetPacket()))

	if err != nil {
		panic(err)
	}

	c.conn.Write([]byte{0})

	fmt.Println("---> NODE INFO CACHED GET:")
	fmt.Println(hex.Dump(serializeBuf.Bytes()))
	fmt.Println("----")

	c.conn.Write(serializeBuf.Bytes())

	fmt.Printf("%d bytes to read", reader.Buffered())

	buf = make([]byte, 64)
	byteCount, err = reader.Read(buf)
	if err != nil {
		fmt.Println("Error reading from client")
		c.Close()
		return
	}

	fmt.Printf("Read %d bytes: %s\n", byteCount, buf)

	// udpPacket := NewUDPv6Packet(44123, 4123, srcAddr, destAddr, zwave.MakeNodeInfoCachedGetPacket())
	// ipPacket := NewIPv6Packet(srcAddr, destAddr, udpPacket.Marshal())
	//
	// c.conn.Write(ipPacket.Marshal())
	// fmt.Println(hex.Dump(ipPacket.Marshal()))

	// for {
	// 	reqLen, err := c.conn.Read(buf)
	// 	if err == io.EOF {
	// 		fmt.Println("Connection closed")
	// 		return
	// 	}
	// 	if err != nil {
	// 		fmt.Println("Connection error:", err.Error())
	// 		return
	// 	}
	//
	// 	if buf[0] == 0x61 && buf[1] == 0x02 && buf[2] == 0xff {
	// 		fmt.Println("---> CONFIGURATION_SET OK")
	// 	} else {
	// 		fmt.Println("---> UNKNOWN PACKET: ", buf[0:reqLen])
	// 	}
	// }
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
