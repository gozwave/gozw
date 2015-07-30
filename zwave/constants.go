package zwave

const (
	FnGetInitAppData                  = 0x02
	FnApplicationNodeInformation      = 0x03
	FnApplicationCommandHandler       = 0x04
	FnSerialApiCapabilities           = 0x07
	FnSendData                        = 0x13
	FnGetVersion                      = 0x15
	FnMemoryGetId                     = 0x20
	FnGetNodeProtocolInfo             = 0x41
	FnSetDefault                      = 0x42
	FnApplicationControllerUpdate     = 0x49
	FnAddNodeToNetwork                = 0x4a
	FnRemoveNodeFromNetwork           = 0x4b
	FnRequestNetworkUpdate            = 0x53
	FnRequestNodeInfo                 = 0x60
	FnRemoveFailingNode               = 0x61
	FnIsNodeFailed                    = 0x62
	FnApplicationCommandHandlerBridge = 0xA8
	FnSerialAPIReady                  = 0xEF
)

const (
	ApplicationNodeInfoNotListening          = 0x00
	ApplicationNodeInfoListening             = 0x01
	ApplicationNodeInfoOptionalFunctionality = 0x02
	ApplicationFreqListeningMode1000ms       = 0x10
	ApplicationFreqListeningMode250ms        = 0x20
)

const (
	UpdateStateNodeInfoReceived  = 0x84
	UpdateStateNodeInfoReqDone   = 0x82
	UpdateStateNodeInfoReqFailed = 0x81
	UpdateStateRoutingPending    = 0x80
	UpdateStateNewIdAssigned     = 0x40
	UpdateStateDeleteDone        = 0x20
	UpdateStateSucId             = 0x10
)

const (
	AddNodeAny        = 1
	AddNodeController = 2
	AddNodeSlave      = 3
	AddNodeExisting   = 4
	AddNodeStop       = 5
	AddNodeStopFailed = 6
)

const (
	TransmitOptionAck       = 0x01
	TransmitOptionLowPower  = 0x02
	TransmitOptionAutoRoute = 0x04
	TransmitOptionNoRoute   = 0x10
	TransmitOptionExplore   = 0x20
)

const (
	AddNodeOptionNormalPower = 0x80
	AddNodeOptionNetworkWide = 0x40
)

const (
	AddNodeStatusLearnReady       = 1
	AddNodeStatusNodeFound        = 2
	AddNodeStatusAddingSlave      = 3
	AddNodeStatusAddingController = 4
	AddNodeStatusProtocolDone     = 5
	AddNodeStatusDone             = 6
	AddNodeStatusFailed           = 7
	AddNodeStatusSecurityFailed   = 9
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

const (
	RemoveNodeStatusLearnReady         = AddNodeStatusLearnReady
	RemoveNodeStatusNodeFound          = AddNodeStatusNodeFound
	RemoveNodeStatusRemovingSlave      = AddNodeStatusAddingSlave
	RemoveNodeStatusRemovingController = AddNodeStatusAddingController
	RemoveNodeStatusProtocolDone       = AddNodeStatusProtocolDone
	RemoveNodeStatusDone               = AddNodeStatusDone
	RemoveNodeStatusFailed             = AddNodeStatusFailed
)
