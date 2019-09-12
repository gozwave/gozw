package serialapi

import (
	"errors"

	"github.com/gozwave/gozw/frame"
	"github.com/gozwave/gozw/protocol"
	"github.com/gozwave/gozw/session"
)

// GetVersion will retrieve version information.
func (s *Layer) GetVersion() (version *Version, err error) {

	s.l.Debug("getting version")

	done := make(chan *frame.Frame)

	request := &session.Request{
		FunctionID: protocol.FnGetVersion,
		HasReturn:  true,
		ReturnCallback: func(err error, ret *frame.Frame) bool {
			done <- ret
			return false
		},
	}

	s.sessionLayer.MakeRequest(request)
	ret := <-done

	if ret == nil {
		return nil, errors.New("Error getting version")
	}

	version = &Version{
		Version:     string(ret.Payload[1:12]),
		LibraryType: ret.Payload[13],
	}

	return
}

// Version contains version information.
type Version struct {
	Version     string
	LibraryType byte
}

// GetLibraryTypeString will return the string representation of the library type.
func (v *Version) GetLibraryTypeString() string {
	switch v.LibraryType {
	case protocol.LibraryControllerStatic:
		return "Static Controller"
	case protocol.LibraryController:
		return "Controller"
	case protocol.LibrarySlaveEnhanced:
		return "Enhanced Slave"
	case protocol.LibrarySlave:
		return "Slave"
	case protocol.LibraryInstaller:
		return "Installer"
	case protocol.LibrarySlaveRouting:
		return "Routing Slave"
	case protocol.LibraryControllerBridge:
		return "Bridge Controller"
	case protocol.LibraryDUT:
		return "DUT"
	case protocol.LibraryAvRemote:
		return "AV Remote"
	case protocol.LibraryAvDevice:
		return "AV Device"
	default:
		return "Unknown"
	}
}
