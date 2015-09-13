// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package thermostatsetpointv3

import "fmt"

type CommandID byte

const (
	CommandGet CommandID = 0x02

	CommandReport CommandID = 0x03

	CommandSet CommandID = 0x01

	CommandSupportedGet CommandID = 0x04

	CommandSupportedReport CommandID = 0x05

	CommandCapabilitiesGet CommandID = 0x09

	CommandCapabilitiesReport CommandID = 0x0A
)

func (c CommandID) String() string {
	switch c {

	case CommandGet:
		return "THERMOSTAT_SETPOINT_GET"

	case CommandReport:
		return "THERMOSTAT_SETPOINT_REPORT"

	case CommandSet:
		return "THERMOSTAT_SETPOINT_SET"

	case CommandSupportedGet:
		return "THERMOSTAT_SETPOINT_SUPPORTED_GET"

	case CommandSupportedReport:
		return "THERMOSTAT_SETPOINT_SUPPORTED_REPORT"

	case CommandCapabilitiesGet:
		return "THERMOSTAT_SETPOINT_CAPABILITIES_GET"

	case CommandCapabilitiesReport:
		return "THERMOSTAT_SETPOINT_CAPABILITIES_REPORT"

	default:
		return fmt.Sprintf("Unknown (0x%X)", byte(c))
	}
}