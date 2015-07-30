package protocol

const (
	FnSerialApiGetInitAppData             uint8 = 0x02
	FnSerialApiApplicationNodeInformation       = 0x03
	FnApplicationCommandHandler                 = 0x04
	FnGetControllerCapabilities                 = 0x05
	FnSerialApiGetTimeouts                      = 0x06
	FnSerialApiGetCapabilities                  = 0x07
	FnSerialApiSoftReset                        = 0x08
	FnGetProtocolVersion                        = 0x09
	FnSendNodeInformation                       = 0x12
	FnSendData                                  = 0x13
	FnSendDataMulti                             = 0x14
	FnGetVersion                                = 0x15
	FnSendDataAbort                             = 0x16
	FnRFPowerLevelSet                           = 0x17
	FnSendDataMeta                              = 0x18
	FnSetRoutingInfo                            = 0x1B
	FnRFPowerLevelRediscoverySet                = 0x1E
	FnMemoryGetId                               = 0x20
	FnGetNodeProtocolInfo                       = 0x41
	FnSetDefault                                = 0x42
	FnAssignReturnRoute                         = 0x46
	FnDeleteReturnRoute                         = 0x47
	FnRequestNodeNeighborUpdate                 = 0x48
	FnApplicationControllerUpdate               = 0x49
	FnAddNodeToNetwork                          = 0x4a
	FnRemoveNodeFromNetwork                     = 0x4b
	FnRequestNetworkUpdate                      = 0x53
	FnRequestNodeInfo                           = 0x60
	FnRemoveFailingNode                         = 0x61
	FnIsNodeFailed                              = 0x62
	FnApplicationCommandHandlerBridge           = 0xA8
	FnSerialAPIReady                            = 0xEF
)

const (
	LibraryControllerStatic uint8 = 0x01
	LibraryController             = 0x02
	LibrarySlaveEnhanced          = 0x03
	LibrarySlave                  = 0x04
	LibraryInstaller              = 0x05
	LibrarySlaveRouting           = 0x06
	LibraryControllerBridge       = 0x07
	LibraryDUT                    = 0x08
	LibraryAvRemote               = 0x0A
	LibraryAvDevice               = 0x0B
)

const (
	AddNodeAny        uint8 = 1
	AddNodeController       = 2
	AddNodeSlave            = 3
	AddNodeExisting         = 4
	AddNodeStop             = 5
	AddNodeStopFailed       = 6
)

const (
	TransmitOptionAck       uint8 = 0x01
	TransmitOptionLowPower        = 0x02
	TransmitOptionAutoRoute       = 0x04
	TransmitOptionNoRoute         = 0x10
	TransmitOptionExplore         = 0x20
)

const (
	TransmitCompleteOk      uint8 = 0x00
	TransmitCompleteNoAck         = 0x01
	TransmitCompleteFail          = 0x02
	TransmitRoutingNotIdle        = 0x03
	TransmitCompleteNoRoute       = 0x04
)

const (
	AddNodeOptionNormalPower uint8 = 0x80
	AddNodeOptionNetworkWide       = 0x40
)

const (
	AddNodeStatusLearnReady       uint8 = 1
	AddNodeStatusNodeFound              = 2
	AddNodeStatusAddingSlave            = 3
	AddNodeStatusAddingController       = 4
	AddNodeStatusProtocolDone           = 5
	AddNodeStatusDone                   = 6
	AddNodeStatusFailed                 = 7
	AddNodeStatusSecurityFailed         = 9
)

const (
	RemoveNodeAny        uint8 = AddNodeAny
	RemoveNodeController       = AddNodeController
	RemoveNodeSlave            = AddNodeSlave
	RemoveNodeStop             = AddNodeStop
)

const (
	RemoveNodeOptionNormalPower uint8 = AddNodeOptionNormalPower
	RemoveNodeOptionNetworkWide       = AddNodeOptionNetworkWide
)

const (
	RemoveNodeStatusLearnReady         uint8 = AddNodeStatusLearnReady
	RemoveNodeStatusNodeFound                = AddNodeStatusNodeFound
	RemoveNodeStatusRemovingSlave            = AddNodeStatusAddingSlave
	RemoveNodeStatusRemovingController       = AddNodeStatusAddingController
	RemoveNodeStatusProtocolDone             = AddNodeStatusProtocolDone
	RemoveNodeStatusDone                     = AddNodeStatusDone
	RemoveNodeStatusFailed                   = AddNodeStatusFailed
)
