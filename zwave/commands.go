package zwave

import "net"

const CommandClassZIPGateway uint8 = 0x61
const GatewayConfigurationSet uint8 = 0x01

func MakeConfigurationSetPacket() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, CommandClassZIPGateway)
	bytes = append(bytes, GatewayConfigurationSet)
	bytes = append(bytes, (make([]byte, 16))...)     // LAN IP Address
	bytes = append(bytes, byte(0x0))                 // LAN IP Prefix Prefix Length
	bytes = append(bytes, net.ParseIP("3000::1")...) // Portal IPv6 Prefix
	bytes = append(bytes, byte(0x40))                // Portal IPv6 Prefix length
	bytes = append(bytes, net.ParseIP("3000::1")...) // Gateway IPv6 Address
	bytes = append(bytes, (make([]byte, 16))...)     // PAN IPv6 Address

	return bytes
}

const CommandClassNetworkManagementProxy uint8 = 0x52
const NodeInfoCachedGet uint8 = 0x03

func MakeNodeInfoCachedGetPacket() []byte {
	bytes := make([]byte, 0)

	bytes = append(bytes, CommandClassNetworkManagementProxy)
	bytes = append(bytes, NodeInfoCachedGet)
	bytes = append(bytes, uint8(0x20))
	bytes = append(bytes, uint8(0x02))
	bytes = append(bytes, uint8(0x00))

	return bytes
}
