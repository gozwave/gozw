package functions

const (
	ZwGetInitData           = 0x02
	ZwAppNodeInfo           = 0x03
	ZwSendData              = 0x13
	ZwGetNodeProtocolInfo   = 0x41
	ZwAddNodeToNetwork      = 0x4a
	ZwRemoveNodeFromNetwork = 0x4b
	ZwRequestNodeInfo       = 0x60
	ZwSerialAPIReady        = 0xEF
)

const (
	AddNodeAny                  = 1
	AddNodeController           = 2
	AddNodeSlave                = 3
	AddNodeExisting             = 4
	AddNodeStop                 = 5
	AddNodeStopFailed           = 6
	AddNodeStatusSecurityFailed = 9
)

const (
	AddNodeOptionNormalPower = 0x80
	AddNodeOptionNetworkWide = 0x40
)

const (
	RemoveNodeAny        = AddNodeAny
	RemoveNodeController = AddNodeController
	RemoveNodeSlave      = AddNodeSlave
	RemoveNodeStop       = AddNodeStop
)

const (
	RemoveNodeOptionNormalPower = AddNodeOptionNormalPower
	RemoveNodeOptionNetworkWide = AddNodeOptionNetworkWide
)
