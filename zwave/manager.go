package zwave

type Manager struct {
	session *SessionLayer

	Version                 byte
	APIType                 string
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
	m.APIType = appInfo.GetAPIType()
	m.TimerFunctionsSupported = appInfo.TimerFunctionsSupported()
	m.IsPrimaryController = appInfo.IsPrimaryController()
	m.NodeList = appInfo.GetNodeIds()
}

func (m *Manager) GetAppInfo() *NodeListResponse {
	resp := m.session.ExecuteCommand(0x02, []byte{})
	respPayload := ParseFunctionPayload(resp.Payload).(*NodeListResponse)

	return respPayload
}
