// THIS FILE IS AUTO-GENERATED
// DO NOT MODIFY

package cc

import "fmt"

const (
	Alarm                                    CommandClassID = 0x71
	AlarmV2                                  CommandClassID = 0x71
	NotificationV3                           CommandClassID = 0x71
	NotificationV4                           CommandClassID = 0x71
	NotificationV5                           CommandClassID = 0x71
	NotificationV6                           CommandClassID = 0x71
	NotificationV7                           CommandClassID = 0x71
	NotificationV8                           CommandClassID = 0x71
	ApplicationStatus                        CommandClassID = 0x22
	AssociationCommandConfiguration          CommandClassID = 0x9B
	Association                              CommandClassID = 0x85
	AssociationV2                            CommandClassID = 0x85
	AvContentDirectoryMd                     CommandClassID = 0x95
	AvContentSearchMd                        CommandClassID = 0x97
	AvRendererStatus                         CommandClassID = 0x96
	AvTaggingMd                              CommandClassID = 0x99
	BasicTariffInfo                          CommandClassID = 0x36
	BasicWindowCovering                      CommandClassID = 0x50
	Basic                                    CommandClassID = 0x20
	BasicV2                                  CommandClassID = 0x20
	Battery                                  CommandClassID = 0x80
	ChimneyFan                               CommandClassID = 0x2A
	ClimateControlSchedule                   CommandClassID = 0x46
	Clock                                    CommandClassID = 0x81
	Configuration                            CommandClassID = 0x70
	ConfigurationV2                          CommandClassID = 0x70
	ConfigurationV3                          CommandClassID = 0x70
	ConfigurationV4                          CommandClassID = 0x70
	ControllerReplication                    CommandClassID = 0x21
	Crc16Encap                               CommandClassID = 0x56
	DcpConfig                                CommandClassID = 0x3A
	DcpMonitor                               CommandClassID = 0x3B
	DoorLockLogging                          CommandClassID = 0x4C
	DoorLock                                 CommandClassID = 0x62
	DoorLockV2                               CommandClassID = 0x62
	DoorLockV3                               CommandClassID = 0x62
	EnergyProduction                         CommandClassID = 0x90
	FirmwareUpdateMd                         CommandClassID = 0x7A
	FirmwareUpdateMdV2                       CommandClassID = 0x7A
	FirmwareUpdateMdV3                       CommandClassID = 0x7A
	FirmwareUpdateMdV4                       CommandClassID = 0x7A
	FirmwareUpdateMdV5                       CommandClassID = 0x7A
	GeographicLocation                       CommandClassID = 0x8C
	GroupingName                             CommandClassID = 0x7B
	Hail                                     CommandClassID = 0x82
	HrvControl                               CommandClassID = 0x39
	HrvStatus                                CommandClassID = 0x37
	Indicator                                CommandClassID = 0x87
	IndicatorV2                              CommandClassID = 0x87
	IpConfiguration                          CommandClassID = 0x9A
	Language                                 CommandClassID = 0x89
	Lock                                     CommandClassID = 0x76
	ManufacturerProprietary                  CommandClassID = 0x91
	ManufacturerSpecific                     CommandClassID = 0x72
	ManufacturerSpecificV2                   CommandClassID = 0x72
	Mark                                     CommandClassID = 0xEF
	MeterPulse                               CommandClassID = 0x35
	MeterTblConfig                           CommandClassID = 0x3C
	MeterTblMonitor                          CommandClassID = 0x3D
	MeterTblMonitorV2                        CommandClassID = 0x3D
	MeterTblPush                             CommandClassID = 0x3E
	Meter                                    CommandClassID = 0x32
	MeterV2                                  CommandClassID = 0x32
	MeterV3                                  CommandClassID = 0x32
	MeterV4                                  CommandClassID = 0x32
	MtpWindowCovering                        CommandClassID = 0x51
	MultiChannelAssociationV2                CommandClassID = 0x8E
	MultiChannelAssociationV3                CommandClassID = 0x8E
	MultiChannelV2                           CommandClassID = 0x60
	MultiChannelV3                           CommandClassID = 0x60
	MultiChannelV4                           CommandClassID = 0x60
	MultiCmd                                 CommandClassID = 0x8F
	MultiInstanceAssociation                 CommandClassID = 0x8E
	MultiInstance                            CommandClassID = 0x60
	NetworkManagementProxy                   CommandClassID = 0x52
	NetworkManagementProxyV2                 CommandClassID = 0x52
	NetworkManagementBasic                   CommandClassID = 0x4D
	NetworkManagementBasicV2                 CommandClassID = 0x4D
	NetworkManagementInclusion               CommandClassID = 0x34
	NetworkManagementInclusionV2             CommandClassID = 0x34
	NoOperation                              CommandClassID = 0x00
	NodeNaming                               CommandClassID = 0x77
	NonInteroperable                         CommandClassID = 0xF0
	Powerlevel                               CommandClassID = 0x73
	PrepaymentEncapsulation                  CommandClassID = 0x41
	Prepayment                               CommandClassID = 0x3F
	Proprietary                              CommandClassID = 0x88
	Protection                               CommandClassID = 0x75
	ProtectionV2                             CommandClassID = 0x75
	RateTblConfig                            CommandClassID = 0x48
	RateTblMonitor                           CommandClassID = 0x49
	RemoteAssociationActivate                CommandClassID = 0x7C
	RemoteAssociation                        CommandClassID = 0x7D
	SceneActivation                          CommandClassID = 0x2B
	SceneActuatorConf                        CommandClassID = 0x2C
	SceneControllerConf                      CommandClassID = 0x2D
	ScheduleEntryLock                        CommandClassID = 0x4E
	ScheduleEntryLockV2                      CommandClassID = 0x4E
	ScheduleEntryLockV3                      CommandClassID = 0x4E
	ScreenAttributes                         CommandClassID = 0x93
	ScreenAttributesV2                       CommandClassID = 0x93
	ScreenMd                                 CommandClassID = 0x92
	ScreenMdV2                               CommandClassID = 0x92
	SecurityPanelMode                        CommandClassID = 0x24
	SecurityPanelZoneSensor                  CommandClassID = 0x2F
	SecurityPanelZone                        CommandClassID = 0x2E
	Security                                 CommandClassID = 0x98
	SensorAlarm                              CommandClassID = 0x9C
	SensorBinary                             CommandClassID = 0x30
	SensorBinaryV2                           CommandClassID = 0x30
	SensorConfiguration                      CommandClassID = 0x9E
	SensorMultilevel                         CommandClassID = 0x31
	SensorMultilevelV2                       CommandClassID = 0x31
	SensorMultilevelV3                       CommandClassID = 0x31
	SensorMultilevelV4                       CommandClassID = 0x31
	SensorMultilevelV5                       CommandClassID = 0x31
	SensorMultilevelV6                       CommandClassID = 0x31
	SensorMultilevelV7                       CommandClassID = 0x31
	SensorMultilevelV8                       CommandClassID = 0x31
	SensorMultilevelV9                       CommandClassID = 0x31
	SensorMultilevelV10                      CommandClassID = 0x31
	SilenceAlarm                             CommandClassID = 0x9D
	SimpleAvControl                          CommandClassID = 0x94
	SwitchAll                                CommandClassID = 0x27
	SwitchBinary                             CommandClassID = 0x25
	SwitchBinaryV2                           CommandClassID = 0x25
	SwitchMultilevel                         CommandClassID = 0x26
	SwitchMultilevelV2                       CommandClassID = 0x26
	SwitchMultilevelV3                       CommandClassID = 0x26
	SwitchMultilevelV4                       CommandClassID = 0x26
	SwitchToggleBinary                       CommandClassID = 0x28
	SwitchToggleMultilevel                   CommandClassID = 0x29
	TariffConfig                             CommandClassID = 0x4A
	TariffTblMonitor                         CommandClassID = 0x4B
	ThermostatFanMode                        CommandClassID = 0x44
	ThermostatFanModeV2                      CommandClassID = 0x44
	ThermostatFanModeV3                      CommandClassID = 0x44
	ThermostatFanModeV4                      CommandClassID = 0x44
	ThermostatFanState                       CommandClassID = 0x45
	ThermostatFanStateV2                     CommandClassID = 0x45
	ThermostatHeating                        CommandClassID = 0x38
	ThermostatMode                           CommandClassID = 0x40
	ThermostatModeV2                         CommandClassID = 0x40
	ThermostatModeV3                         CommandClassID = 0x40
	ThermostatOperatingState                 CommandClassID = 0x42
	ThermostatOperatingStateV2               CommandClassID = 0x42
	ThermostatSetback                        CommandClassID = 0x47
	ThermostatSetpoint                       CommandClassID = 0x43
	ThermostatSetpointV2                     CommandClassID = 0x43
	ThermostatSetpointV3                     CommandClassID = 0x43
	TimeParameters                           CommandClassID = 0x8B
	Time                                     CommandClassID = 0x8A
	TimeV2                                   CommandClassID = 0x8A
	TransportServiceV2                       CommandClassID = 0x55
	TransportService                         CommandClassID = 0x55
	UserCode                                 CommandClassID = 0x63
	Version                                  CommandClassID = 0x86
	VersionV2                                CommandClassID = 0x86
	WakeUp                                   CommandClassID = 0x84
	WakeUpV2                                 CommandClassID = 0x84
	ZensorNet                                CommandClassID = 0x02
	Zip6lowpan                               CommandClassID = 0x4F
	Zip                                      CommandClassID = 0x23
	ZipV2                                    CommandClassID = 0x23
	ZipV3                                    CommandClassID = 0x23
	ZwaveCmdClass                            CommandClassID = 0x01
	ApplicationCapability                    CommandClassID = 0x57
	SwitchColor                              CommandClassID = 0x33
	SwitchColorV2                            CommandClassID = 0x33
	SwitchColorV3                            CommandClassID = 0x33
	Schedule                                 CommandClassID = 0x53
	ScheduleV2                               CommandClassID = 0x53
	ScheduleV3                               CommandClassID = 0x53
	NetworkManagementPrimary                 CommandClassID = 0x54
	ZipNd                                    CommandClassID = 0x58
	AssociationGrpInfo                       CommandClassID = 0x59
	AssociationGrpInfoV2                     CommandClassID = 0x59
	AssociationGrpInfoV3                     CommandClassID = 0x59
	DeviceResetLocally                       CommandClassID = 0x5A
	CentralScene                             CommandClassID = 0x5B
	CentralSceneV2                           CommandClassID = 0x5B
	CentralSceneV3                           CommandClassID = 0x5B
	IpAssociation                            CommandClassID = 0x5C
	Antitheft                                CommandClassID = 0x5D
	AntitheftV2                              CommandClassID = 0x5D
	ZwaveplusInfo                            CommandClassID = 0x5E
	ZwaveplusInfoV2                          CommandClassID = 0x5E
	ZipGateway                               CommandClassID = 0x5F
	ZipPortal                                CommandClassID = 0x61
	Dmx                                      CommandClassID = 0x65
	BarrierOperator                          CommandClassID = 0x66
	NetworkManagementInstallationMaintenance CommandClassID = 0x67
	ZipNaming                                CommandClassID = 0x68
	Mailbox                                  CommandClassID = 0x69
	WindowCovering                           CommandClassID = 0x6A
	Security2                                CommandClassID = 0x9F
	Irrigation                               CommandClassID = 0x6B
	Supervision                              CommandClassID = 0x6C
	HumidityControlSetpoint                  CommandClassID = 0x64
	HumidityControlMode                      CommandClassID = 0x6D
	HumidityControlOperatingState            CommandClassID = 0x6E
	EntryControl                             CommandClassID = 0x6F
	InclusionController                      CommandClassID = 0x74
)

func (c CommandClassID) String() string {
	switch c {

	case Alarm:
		return "Command Class Alarm"

	case ApplicationStatus:
		return "Command Class Application Status"

	case AssociationCommandConfiguration:
		return "Command Class Association Command Configuration"

	case Association:
		return "Command Class Association"

	case AvContentDirectoryMd:
		return "Command Class Av Content Directory Md"

	case AvContentSearchMd:
		return "Command Class Av Content Search Md"

	case AvRendererStatus:
		return "Command Class Av Renderer Status"

	case AvTaggingMd:
		return "Command Class Av Tagging Md"

	case BasicTariffInfo:
		return "Command Class Basic Tariff Info"

	case BasicWindowCovering:
		return "Command Class Basic Window Covering"

	case Basic:
		return "Command Class Basic"

	case Battery:
		return "Command Class Battery"

	case ChimneyFan:
		return "Command Class Chimney Fan"

	case ClimateControlSchedule:
		return "Command Class Climate Control Schedule"

	case Clock:
		return "Command Class Clock"

	case Configuration:
		return "Command Class Configuration"

	case ControllerReplication:
		return "Command Class Controller Replication"

	case Crc16Encap:
		return "Command Class CRC16 Encap"

	case DcpConfig:
		return "Command Class Dcp Config"

	case DcpMonitor:
		return "Command Class Dcp Monitor"

	case DoorLockLogging:
		return "Command Class Door Lock Logging"

	case DoorLock:
		return "Command Class Door Lock"

	case EnergyProduction:
		return "Command Class Energy Production"

	case FirmwareUpdateMd:
		return "Command Class Firmware Update Md"

	case GeographicLocation:
		return "Command Class Geographic Location"

	case GroupingName:
		return "Command Class Grouping Name"

	case Hail:
		return "Command Class Hail"

	case HrvControl:
		return "Command Class Hrv Control"

	case HrvStatus:
		return "Command Class Hrv Status"

	case Indicator:
		return "Command Class Indicator"

	case IpConfiguration:
		return "Command Class Ip Configuration"

	case Language:
		return "Command Class Language"

	case Lock:
		return "Command Class Lock"

	case ManufacturerProprietary:
		return "Command Class Manufacturer Proprietary"

	case ManufacturerSpecific:
		return "Command Class Manufacturer Specific"

	case Mark:
		return "Command Class Mark"

	case MeterPulse:
		return "Command Class Meter Pulse"

	case MeterTblConfig:
		return "Command Class Meter Tbl Config"

	case MeterTblMonitor:
		return "Command Class Meter Tbl Monitor"

	case MeterTblPush:
		return "Command Class Meter Tbl Push"

	case Meter:
		return "Command Class Meter"

	case MtpWindowCovering:
		return "Command Class Mtp Window Covering"

	case MultiCmd:
		return "Command Class Multi Cmd"

	case MultiInstanceAssociation:
		return "Command Class Multi Instance Association"

	case MultiInstance:
		return "Command Class Multi Instance"

	case NetworkManagementProxy:
		return "Command Class Network Management Proxy"

	case NetworkManagementBasic:
		return "Command Class Network Management Basic"

	case NetworkManagementInclusion:
		return "Command Class Network Management Inclusion"

	case NoOperation:
		return "NOP"

	case NodeNaming:
		return "Command Class Node Naming"

	case NonInteroperable:
		return "Command Class Non Interoperable"

	case Powerlevel:
		return "Command Class Powerlevel"

	case PrepaymentEncapsulation:
		return "Command Class Prepayment Encapsulation"

	case Prepayment:
		return "Command Class Prepayment"

	case Proprietary:
		return "Command Class Proprietary"

	case Protection:
		return "Command Class Protection"

	case RateTblConfig:
		return "Command Class Rate Tbl Config"

	case RateTblMonitor:
		return "Command Class Rate Tbl Monitor"

	case RemoteAssociationActivate:
		return "Command Class Remote Association Activate"

	case RemoteAssociation:
		return "Command Class Remote Association"

	case SceneActivation:
		return "Command Class Scene Activation"

	case SceneActuatorConf:
		return "Command Class Scene Actuator Conf"

	case SceneControllerConf:
		return "Command Class Scene Controller Conf"

	case ScheduleEntryLock:
		return "Command Class Schedule Entry Lock"

	case ScreenAttributes:
		return "Command Class Screen Attributes"

	case ScreenMd:
		return "Command Class Screen Md"

	case SecurityPanelMode:
		return "Command Class Security Panel Mode"

	case SecurityPanelZoneSensor:
		return "Command Class Security Panel Zone Sensor"

	case SecurityPanelZone:
		return "Command Class Security Panel Zone"

	case Security:
		return "Command Class Security"

	case SensorAlarm:
		return "Command Class Sensor Alarm"

	case SensorBinary:
		return "Command Class Sensor Binary"

	case SensorConfiguration:
		return "Command Class Sensor Configuration"

	case SensorMultilevel:
		return "Command Class Sensor Multilevel"

	case SilenceAlarm:
		return "Command Class Silence Alarm"

	case SimpleAvControl:
		return "Command Class Simple Av Control"

	case SwitchAll:
		return "Command Class Switch All"

	case SwitchBinary:
		return "Command Class Switch Binary"

	case SwitchMultilevel:
		return "Command Class Switch Multilevel"

	case SwitchToggleBinary:
		return "Command Class Switch Toggle Binary"

	case SwitchToggleMultilevel:
		return "Command Class Switch Toggle Multilevel"

	case TariffConfig:
		return "Command Class Tariff Config"

	case TariffTblMonitor:
		return "Command Class Tariff Tbl Monitor"

	case ThermostatFanMode:
		return "Command Class Thermostat Fan Mode"

	case ThermostatFanState:
		return "Command Class Thermostat Fan State"

	case ThermostatHeating:
		return "Command Class Thermostat Heating"

	case ThermostatMode:
		return "Command Class Thermostat Mode"

	case ThermostatOperatingState:
		return "Command Class Thermostat Operating State"

	case ThermostatSetback:
		return "Command Class Thermostat Setback"

	case ThermostatSetpoint:
		return "Command Class Thermostat Setpoint"

	case TimeParameters:
		return "Command Class Time Parameters"

	case Time:
		return "Command Class Time"

	case TransportService:
		return "Command Class Transport Service"

	case UserCode:
		return "Command Class User Code"

	case Version:
		return "Command Class Version"

	case WakeUp:
		return "Command Class Wake Up"

	case ZensorNet:
		return "Command Class Zensor Net"

	case Zip6lowpan:
		return "Command Class Z/IP 6lowpan"

	case Zip:
		return "Command Class Z/IP"

	case ZwaveCmdClass:
		return "Z-Wave protocol Command Class"

	case ApplicationCapability:
		return "Command Class Application Capability"

	case SwitchColor:
		return "Command Class Switch Color"

	case Schedule:
		return "Command Class Schedule"

	case NetworkManagementPrimary:
		return "Command Class Network Management Primary"

	case ZipNd:
		return "Command Class Z/IP-ND"

	case AssociationGrpInfo:
		return "Command Class Association Group Info"

	case DeviceResetLocally:
		return "Command Class Device Reset Locally"

	case CentralScene:
		return "Command Class Central Scene"

	case IpAssociation:
		return "Command Class Ip Association"

	case Antitheft:
		return "Command Class Anti-theft"

	case ZwaveplusInfo:
		return "Command Class Z-Wave+ Info"

	case ZipGateway:
		return "Command Class Z/IP Gateway"

	case ZipPortal:
		return "Command Class Z/IP Portal"

	case Dmx:
		return "Command Class DMX"

	case BarrierOperator:
		return "Command Class Barrier Operator"

	case NetworkManagementInstallationMaintenance:
		return "Command Class Network Management Installation and Maintenance"

	case ZipNaming:
		return "Command Class Z/IP Naming and Location"

	case Mailbox:
		return "Command Class Mailbox"

	case WindowCovering:
		return "Command Class Window Covering"

	case Security2:
		return "Command Class Security 2"

	case Irrigation:
		return "Command Class Irrigation"

	case Supervision:
		return "Command Class Supervision"

	case HumidityControlSetpoint:
		return "Command Class Humidity Control Setpoint"

	case HumidityControlMode:
		return "Command Class Humidity Control Mode"

	case HumidityControlOperatingState:
		return "Command Class Humidity Control Operating State"

	case EntryControl:
		return "Command Class Entry Control"

	case InclusionController:
		return "Command Class Inclusion Controller"

	default:
		return fmt.Sprintf("Unknown (0x%X)", byte(c))
	}
}
