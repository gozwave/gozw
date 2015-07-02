package zwave

import "fmt"

type Manager struct {
	session *SessionLayer

	ApiVersion     string
	ApiLibraryType string

	HomeId uint32
	NodeId byte

	Version                 byte
	ApiType                 string
	TimerFunctionsSupported bool
	IsPrimaryController     bool
	NodeList                []uint8
	ApplicationVersion      byte
	ApplicationRevision     byte
	SupportedFunctions      []byte
}

func NewManager(session *SessionLayer) *Manager {
	manager := &Manager{
		session: session,
	}

	manager.Init()

	return manager
}

func (m *Manager) Init() {
	version := m.GetVersion()
	m.ApiVersion = version.ApiVersion
	m.ApiLibraryType = version.GetLibraryTypeString()

	ids := m.GetHomeId()
	m.HomeId = ids.HomeId
	m.NodeId = ids.NodeId

	appInfo := m.GetAppInfo()
	m.Version = appInfo.Version
	m.ApiType = appInfo.GetApiType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.NodeList = appInfo.GetNodeIds()

	serialApi := m.GetSerialApiCapabilities()
	m.ApplicationVersion = serialApi.ApplicationVersion
	m.ApplicationRevision = serialApi.ApplicationRevision
	m.SupportedFunctions = serialApi.GetSupportedFunctions()
}

func (m *Manager) SetApplicationNodeInformation() {
	m.session.ExecuteCommandNoWait(FnApplicationNodeInformation, []byte{
		ApplicationNodeInfoListening | ApplicationFreqListeningMode250ms | ApplicationNodeInfoOptionalFunctionality,
		GenericTypeGenericController,
		SpecificTypePortableSceneController,
		0x01,
		CommandClassBasic,
	})
}

func (m *Manager) FactoryReset() {
	m.session.ExecuteCommand(FnSetDefault, []byte{0x01})
}

func (m *Manager) AddNode() {
	m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
		AddNodeAny | AddNodeOptionNetworkWide | AddNodeOptionNormalPower,
		0x01,
	})

	for {
		frame := m.session.WaitForFrame()
		callback := ParseFunctionPayload(frame.Payload).(*AddRemoveNodeCallback)

		switch {
		case callback.Status == AddNodeStatusLearnReady:
			fmt.Print("Add node ready... ")
		case callback.Status == AddNodeStatusNodeFound:
			fmt.Print("found node... ")
		case callback.Status == AddNodeStatusAddingSlave:
			fmt.Print("slave... ")
		case callback.Status == AddNodeStatusProtocolDone:
			fmt.Println("protocol done")
			m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
				AddNodeStop,
				0x02,
			})
		case callback.Status == AddNodeStatusDone:
			fmt.Println("done")
			m.session.ExecuteCommandNoWait(FnAddNodeToNetwork, []byte{
				AddNodeStop,
				0x0,
			})

			return
		default:
			fmt.Println("unknown frame", callback)
		}

	}
}

func (m *Manager) RemoveNode() {
	m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
		RemoveNodeAny | RemoveNodeOptionNetworkWide | RemoveNodeOptionNormalPower,
		0x01,
	})

	for {
		frame := m.session.WaitForFrame()
		callback := ParseFunctionPayload(frame.Payload).(*AddRemoveNodeCallback)

		switch {
		case callback.Status == RemoveNodeStatusLearnReady:
			fmt.Print("Remove node ready... ")
		case callback.Status == RemoveNodeStatusNodeFound:
			fmt.Print("found node... ")
		case callback.Status == RemoveNodeStatusRemovingSlave:
			fmt.Print("slave... ")
		case callback.Status == RemoveNodeStatusProtocolDone:
			fmt.Println("protocol done")
			m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
				RemoveNodeStop,
				0x02,
			})
		case callback.Status == RemoveNodeStatusDone:
			fmt.Println("done")
			m.session.ExecuteCommandNoWait(FnRemoveNodeFromNetwork, []byte{
				RemoveNodeStop,
				0x0,
			})

			return
		default:
			fmt.Println("unknown frame", callback)
		}

	}
}
func (m *Manager) GetHomeId() *MemoryGetIdResponse {
	resp := m.session.ExecuteCommand(FnMemoryGetId, []byte{})
	return ParseFunctionPayload(resp.Payload).(*MemoryGetIdResponse)
}

func (m *Manager) GetAppInfo() *NodeListResponse {
	resp := m.session.ExecuteCommand(FnGetInitAppData, []byte{})
	return ParseFunctionPayload(resp.Payload).(*NodeListResponse)
}

func (m *Manager) GetSerialApiCapabilities() *SerialApiCapabilitiesResponse {
	resp := m.session.ExecuteCommand(FnSerialApiCapabilities, []byte{})
	return ParseFunctionPayload(resp.Payload).(*SerialApiCapabilitiesResponse)
}

func (m *Manager) GetVersion() *VersionResponse {
	resp := m.session.ExecuteCommand(FnGetVersion, []byte{})
	return ParseFunctionPayload(resp.Payload).(*VersionResponse)
}

func (m *Manager) GetNodeProtocolInfo(nodeId uint8) *NodeProtocolInfoResponse {
	resp := m.session.ExecuteCommand(FnGetNodeProtocolInfo, []byte{nodeId})
	return ParseFunctionPayload(resp.Payload).(*NodeProtocolInfoResponse)
}
