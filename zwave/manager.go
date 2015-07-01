package zwave

import "fmt"

type Manager struct {
	session *SessionLayer

	Version                 byte
	ApiType                 string
	TimerFunctionsSupported bool
	IsPrimaryController     bool
	NodeList                []int
}

func NewManager(session *SessionLayer) *Manager {
	manager := &Manager{
		session: session,
	}

	manager.Init()

	return manager
}

func (m *Manager) Init() {
	appInfo := m.GetAppInfo()
	m.Version = appInfo.Version
	m.ApiType = appInfo.GetApiType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.NodeList = appInfo.GetNodeIds()

	serialApi := m.GetSerialApiCapabilities()
	fmt.Println(serialApi)
}

func (m *Manager) GetAppInfo() *NodeListResponse {
	resp := m.session.ExecuteCommand(0x02, []byte{})
	respPayload := ParseFunctionPayload(resp.Payload).(*NodeListResponse)

	return respPayload
}

func (m *Manager) GetSerialApiCapabilities() *SerialApiCapabilitiesResponse {
	resp := m.session.ExecuteCommand(0x07, []byte{})
	respPayload := ParseFunctionPayload(resp.Payload).(*SerialApiCapabilitiesResponse)

	return respPayload
}
