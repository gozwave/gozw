package zwave

const (
	FrameSOFData uint8 = 0x01
	FrameSOFAck  uint8 = 0x06
	FrameSOFNak  uint8 = 0x15
	FrameSOFCan  uint8 = 0x18
)

const (
	FrameTypeReq uint8 = 0x00
	FrameTypeRes uint8 = 0x01
)

const (
	FnGetInitAppData        = 0x02
	FnAppNodeInfo           = 0x03
	FnSerialApiCapabilities = 0x07
	FnSendData              = 0x13
	FnGetNodeProtocolInfo   = 0x41
	FnAddNodeToNetwork      = 0x4a
	FnRemoveNodeFromNetwork = 0x4b
	FnRequestNetworkUpdate  = 0x53
	FnRequestNodeInfo       = 0x60
	FnRemoveFailingNode     = 0x61
	FnIsNodeFailed          = 0x62
	FnSerialAPIReady        = 0xEF
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
