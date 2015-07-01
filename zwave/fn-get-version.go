package zwave

type VersionResponse struct {
	CommandId   byte
	ApiVersion  string
	LibraryType byte
}

func ParseVersionResponse(payload []byte) *VersionResponse {
	val := &VersionResponse{
		CommandId:   payload[0],
		ApiVersion:  string(payload[1:12]),
		LibraryType: payload[13],
	}

	return val
}

func (v *VersionResponse) GetLibraryTypeString() string {
	switch v.LibraryType {
	case LibraryControllerStatic:
		return "Static Controller"
	case LibraryController:
		return "Controller"
	case LibrarySlaveEnhanced:
		return "Enhanced Slave"
	case LibrarySlave:
		return "Slave"
	case LibraryInstaller:
		return "Installer"
	case LibrarySlaveRouting:
		return "Routing Slave"
	case LibraryControllerBridge:
		return "Bridge Controller"
	case LibraryDUT:
		return "DUT"
	case LibraryAvRemote:
		return "AV Remote"
	case LibraryAvDevice:
		return "AV Device"
	default:
		return "Unknown"
	}
}
