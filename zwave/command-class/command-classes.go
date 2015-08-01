package commandclass

import "fmt"

const (
	CommandClassNoOperation                              byte = 0x00
	CommandClassBasic                                         = 0x20
	CommandClassBasicV2                                       = 0x20
	CommandClassControllerReplication                         = 0x21
	CommandClassApplicationStatus                             = 0x22
	CommandClassZip                                           = 0x23 /*OBSOLETED*/
	CommandClassZipV2                                         = 0x23
	CommandClassSecurityPanelMode                             = 0x24
	CommandClassSwitchBinary                                  = 0x25
	CommandClassSwitchBinaryV2                                = 0x25
	CommandClassSwitchMultilevel                              = 0x26
	CommandClassSwitchMultilevelV2                            = 0x26
	CommandClassSwitchMultilevelV3                            = 0x26
	CommandClassSwitchMultilevelV4                            = 0x26
	CommandClassSwitchAll                                     = 0x27
	CommandClassSwitchToggleBinary                            = 0x28 /*DEPRECATED*/
	CommandClassSwitchToggleMultilevel                        = 0x29 /*DEPRECATED*/
	CommandClassChimneyFan                                    = 0x2A
	CommandClassSceneActivation                               = 0x2B
	CommandClassSceneActuatorConf                             = 0x2C
	CommandClassSceneControllerConf                           = 0x2D
	CommandClassSecurityPanelZone                             = 0x2E
	CommandClassSecurityPanelZoneSensor                       = 0x2F
	CommandClassSensorBinary                                  = 0x30 /*DEPRECATED*/
	CommandClassSensorBinaryV2                                = 0x30 /*DEPRECATED*/
	CommandClassSensorMultilevel                              = 0x31
	CommandClassSensorMultilevelV2                            = 0x31
	CommandClassSensorMultilevelV3                            = 0x31
	CommandClassSensorMultilevelV4                            = 0x31
	CommandClassSensorMultilevelV5                            = 0x31
	CommandClassSensorMultilevelV6                            = 0x31
	CommandClassSensorMultilevelV7                            = 0x31
	CommandClassMeter                                         = 0x32
	CommandClassMeterV2                                       = 0x32
	CommandClassMeterV3                                       = 0x32
	CommandClassMeterV4                                       = 0x32
	CommandClassSwitchColor                                   = 0x33
	CommandClassSwitchColorV2                                 = 0x33
	CommandClassSwitchColorV3                                 = 0x33
	CommandClassNetworkManagementInclusion                    = 0x34
	CommandClassMeterPulse                                    = 0x35 /*DEPRECATED*/
	CommandClassBasicTariffInfo                               = 0x36
	CommandClassHrvStatus                                     = 0x37
	CommandClassThermostatHeating                             = 0x38
	CommandClassHrvControl                                    = 0x39
	CommandClassDcpConfig                                     = 0x3A
	CommandClassDcpMonitor                                    = 0x3B
	CommandClassMeterTblConfig                                = 0x3C
	CommandClassMeterTblMonitor                               = 0x3D
	CommandClassMeterTblMonitorV2                             = 0x3D
	CommandClassMeterTblPush                                  = 0x3E
	CommandClassPrepayment                                    = 0x3F
	CommandClassThermostatMode                                = 0x40
	CommandClassThermostatModeV2                              = 0x40
	CommandClassThermostatModeV3                              = 0x40
	CommandClassPrepaymentEncapsulation                       = 0x41
	CommandClassThermostatOperatingState                      = 0x42
	CommandClassThermostatOperatingStateV2                    = 0x42
	CommandClassThermostatSetpoint                            = 0x43
	CommandClassThermostatSetpointV2                          = 0x43
	CommandClassThermostatSetpointV3                          = 0x43
	CommandClassThermostatFanMode                             = 0x44
	CommandClassThermostatFanModeV2                           = 0x44
	CommandClassThermostatFanModeV3                           = 0x44
	CommandClassThermostatFanModeV4                           = 0x44
	CommandClassThermostatFanState                            = 0x45
	CommandClassThermostatFanStateV2                          = 0x45
	CommandClassClimateControlSchedule                        = 0x46 /*DEPRECATED*/
	CommandClassThermostatSetback                             = 0x47
	CommandClassRateTblConfig                                 = 0x48
	CommandClassRateTblMonitor                                = 0x49
	CommandClassTariffConfig                                  = 0x4A
	CommandClassTariffTblMonitor                              = 0x4B
	CommandClassDoorLockLogging                               = 0x4C
	CommandClassNetworkManagementBasic                        = 0x4D
	CommandClassScheduleEntryLock                             = 0x4E /*DEPRECATED*/
	CommandClassScheduleEntryLockV2                           = 0x4E /*DEPRECATED*/
	CommandClassScheduleEntryLockV3                           = 0x4E /*DEPRECATED*/
	CommandClassZip6Lowpan                                    = 0x4F
	CommandClassBasicWindowCovering                           = 0x50 /*Obsoleted*/
	CommandClassMtpWindowCovering                             = 0x51 /*OBSOLETED*/
	CommandClassNetworkManagementProxy                        = 0x52
	CommandClassSchedule                                      = 0x53
	CommandClassNetworkManagementPrimary                      = 0x54
	CommandClassTransportService                              = 0x55
	CommandClassTransportServiceV2                            = 0x55
	CommandClassCrc16Encap                                    = 0x56
	CommandClassApplicationCapability                         = 0x57
	CommandClassZipNd                                         = 0x58
	CommandClassAssociationGrpInfo                            = 0x59
	CommandClassAssociationGrpInfoV2                          = 0x59
	CommandClassDeviceResetLocally                            = 0x5A
	CommandClassCentralScene                                  = 0x5B
	CommandClassCentralSceneV2                                = 0x5B
	CommandClassIpAssociation                                 = 0x5C
	CommandClassAntitheft                                     = 0x5D
	CommandClassAntitheftV2                                   = 0x5D
	CommandClassZwaveplusInfo                                 = 0x5E /*SDS11907-3*/
	CommandClassZwaveplusInfoV2                               = 0x5E /*SDS11907-3*/
	CommandClassZipGateway                                    = 0x5F
	CommandClassMultiChannelV2                                = 0x60
	CommandClassMultiChannelV3                                = 0x60
	CommandClassMultiChannelV4                                = 0x60
	CommandClassMultiInstance                                 = 0x60 /*Discontinued*/
	CommandClassZipPortal                                     = 0x61
	CommandClassDoorLock                                      = 0x62
	CommandClassDoorLockV2                                    = 0x62
	CommandClassDoorLockV3                                    = 0x62
	CommandClassUserCode                                      = 0x63
	CommandClassAppliance                                     = 0x64
	CommandClassDmx                                           = 0x65
	CommandClassBarrierOperator                               = 0x66
	CommandClassNetworkManagementInstallationMaintenance      = 0x67
	CommandClassZipNaming                                     = 0x68
	CommandClassMailbox                                       = 0x69
	CommandClassWindowCovering                                = 0x6A
	CommandClassConfiguration                                 = 0x70
	CommandClassConfigurationV2                               = 0x70
	CommandClassConfigurationV3                               = 0x70
	CommandClassAlarm                                         = 0x71 /*DEPRECATED*/
	CommandClassAlarmV2                                       = 0x71 /*DEPRECATED*/
	CommandClassNotificationV3                                = 0x71
	CommandClassNotificationV4                                = 0x71
	CommandClassNotificationV5                                = 0x71
	CommandClassManufacturerSpecific                          = 0x72
	CommandClassManufacturerSpecificV2                        = 0x72
	CommandClassPowerlevel                                    = 0x73
	CommandClassProtection                                    = 0x75
	CommandClassProtectionV2                                  = 0x75
	CommandClassLock                                          = 0x76
	CommandClassNodeNaming                                    = 0x77
	CommandClassFirmwareUpdateMd                              = 0x7A /*DEPRECATED*/
	CommandClassFirmwareUpdateMdV2                            = 0x7A
	CommandClassFirmwareUpdateMdV3                            = 0x7A
	CommandClassFirmwareUpdateMdV4                            = 0x7A
	CommandClassGroupingName                                  = 0x7B /*DEPRECATED*/
	CommandClassRemoteAssociationActivate                     = 0x7C /*OBSOLETED*/
	CommandClassRemoteAssociation                             = 0x7D /*OBSOLETED*/
	CommandClassBattery                                       = 0x80
	CommandClassClock                                         = 0x81
	CommandClassHail                                          = 0x82 /*DEPRECATED*/
	CommandClassWakeUp                                        = 0x84
	CommandClassWakeUpV2                                      = 0x84
	CommandClassAssociation                                   = 0x85
	CommandClassAssociationV2                                 = 0x85
	CommandClassVersion                                       = 0x86
	CommandClassVersionV2                                     = 0x86
	CommandClassIndicator                                     = 0x87
	CommandClassProprietary                                   = 0x88 /*DEPRECATED*/
	CommandClassLanguage                                      = 0x89
	CommandClassTime                                          = 0x8A
	CommandClassTimeV2                                        = 0x8A
	CommandClassTimeParameters                                = 0x8B
	CommandClassGeographicLocation                            = 0x8C
	CommandClassMultiChannelAssociationV2                     = 0x8E
	CommandClassMultiChannelAssociationV3                     = 0x8E
	CommandClassMultiInstanceAssociation                      = 0x8E /*Discontinued*/
	CommandClassMultiCmd                                      = 0x8F
	CommandClassEnergyProduction                              = 0x90
	CommandClassManufacturerProprietary                       = 0x91
	CommandClassScreenMd                                      = 0x92
	CommandClassScreenMdV2                                    = 0x92
	CommandClassScreenAttributes                              = 0x93
	CommandClassScreenAttributesV2                            = 0x93
	CommandClassSimpleAvControl                               = 0x94
	CommandClassAvContentDirectoryMd                          = 0x95
	CommandClassAvRendererStatus                              = 0x96
	CommandClassAvContentSearchMd                             = 0x97
	CommandClassSecurity                                      = 0x98
	CommandClassAvTaggingMd                                   = 0x99
	CommandClassIpConfiguration                               = 0x9A
	CommandClassAssociationCommandConfiguration               = 0x9B
	CommandClassSensorAlarm                                   = 0x9C /*SDS10963-4 The Sensor Alarm command class can be used to realize Sensor Alarms.*/
	CommandClassSilenceAlarm                                  = 0x9D /*SDS10963-4 The Alarm Silence command class can be used to nuisance silence to temporarily disable the sounding*/
	CommandClassSensorConfiguration                           = 0x9E /*This command class adds the possibility for sensors to act on either a measured value or on a*/
	CommandClassMark                                          = 0xEF
	CommandClassNonInteroperable                              = 0xF0
)

var CommandClassNames map[byte]string = map[byte]string{
	CommandClassNoOperation:                              "NoOperation",
	CommandClassBasic:                                    "Basic",
	CommandClassControllerReplication:                    "ControllerReplication",
	CommandClassApplicationStatus:                        "ApplicationStatus",
	CommandClassZip:                                      "Zip",
	CommandClassSecurityPanelMode:                        "SecurityPanelMode",
	CommandClassSwitchBinary:                             "SwitchBinary",
	CommandClassSwitchMultilevel:                         "SwitchMultilevel",
	CommandClassSwitchAll:                                "SwitchAll",
	CommandClassSwitchToggleBinary:                       "SwitchToggleBinary",
	CommandClassSwitchToggleMultilevel:                   "SwitchToggleMultilevel",
	CommandClassChimneyFan:                               "ChimneyFan",
	CommandClassSceneActivation:                          "SceneActivation",
	CommandClassSceneActuatorConf:                        "SceneActuatorConf",
	CommandClassSceneControllerConf:                      "SceneControllerConf",
	CommandClassSecurityPanelZone:                        "SecurityPanelZone",
	CommandClassSecurityPanelZoneSensor:                  "SecurityPanelZoneSensor",
	CommandClassSensorBinary:                             "SensorBinary",
	CommandClassSensorMultilevel:                         "SensorMultilevel",
	CommandClassMeter:                                    "Meter",
	CommandClassSwitchColor:                              "SwitchColor",
	CommandClassNetworkManagementInclusion:               "NetworkManagementInclusion",
	CommandClassMeterPulse:                               "MeterPulse",
	CommandClassBasicTariffInfo:                          "BasicTariffInfo",
	CommandClassHrvStatus:                                "HrvStatus",
	CommandClassThermostatHeating:                        "ThermostatHeating",
	CommandClassHrvControl:                               "HrvControl",
	CommandClassDcpConfig:                                "DcpConfig",
	CommandClassDcpMonitor:                               "DcpMonitor",
	CommandClassMeterTblConfig:                           "MeterTblConfig",
	CommandClassMeterTblMonitor:                          "MeterTblMonitor",
	CommandClassMeterTblPush:                             "MeterTblPush",
	CommandClassPrepayment:                               "Prepayment",
	CommandClassThermostatMode:                           "ThermostatMode",
	CommandClassPrepaymentEncapsulation:                  "PrepaymentEncapsulation",
	CommandClassThermostatOperatingState:                 "ThermostatOperatingState",
	CommandClassThermostatSetpoint:                       "ThermostatSetpoint",
	CommandClassThermostatFanMode:                        "ThermostatFanMode",
	CommandClassThermostatFanState:                       "ThermostatFanState",
	CommandClassClimateControlSchedule:                   "ClimateControlSchedule",
	CommandClassThermostatSetback:                        "ThermostatSetback",
	CommandClassRateTblConfig:                            "RateTblConfig",
	CommandClassRateTblMonitor:                           "RateTblMonitor",
	CommandClassTariffConfig:                             "TariffConfig",
	CommandClassTariffTblMonitor:                         "TariffTblMonitor",
	CommandClassDoorLockLogging:                          "DoorLockLogging",
	CommandClassNetworkManagementBasic:                   "NetworkManagementBasic",
	CommandClassScheduleEntryLock:                        "ScheduleEntryLock",
	CommandClassZip6Lowpan:                               "Zip6Lowpan",
	CommandClassBasicWindowCovering:                      "BasicWindowCovering",
	CommandClassMtpWindowCovering:                        "MtpWindowCovering",
	CommandClassNetworkManagementProxy:                   "NetworkManagementProxy",
	CommandClassSchedule:                                 "Schedule",
	CommandClassNetworkManagementPrimary:                 "NetworkManagementPrimary",
	CommandClassTransportService:                         "TransportService",
	CommandClassCrc16Encap:                               "Crc16Encap",
	CommandClassApplicationCapability:                    "ApplicationCapability",
	CommandClassZipNd:                                    "ZipNd",
	CommandClassAssociationGrpInfo:                       "AssociationGrpInfo",
	CommandClassDeviceResetLocally:                       "DeviceResetLocally",
	CommandClassCentralScene:                             "CentralScene",
	CommandClassIpAssociation:                            "IpAssociation",
	CommandClassAntitheft:                                "Antitheft",
	CommandClassZwaveplusInfo:                            "ZwaveplusInfo",
	CommandClassZipGateway:                               "ZipGateway",
	CommandClassMultiChannelV2:                           "MultiChannelV2",
	CommandClassZipPortal:                                "ZipPortal",
	CommandClassDoorLock:                                 "DoorLock",
	CommandClassUserCode:                                 "UserCode",
	CommandClassAppliance:                                "Appliance",
	CommandClassDmx:                                      "Dmx",
	CommandClassBarrierOperator:                          "BarrierOperator",
	CommandClassNetworkManagementInstallationMaintenance: "NetworkManagementInstallationMaintenance",
	CommandClassZipNaming:                                "ZipNaming",
	CommandClassMailbox:                                  "Mailbox",
	CommandClassWindowCovering:                           "WindowCovering",
	CommandClassConfiguration:                            "Configuration",
	CommandClassNotificationV3:                           "Notification",
	CommandClassManufacturerSpecific:                     "ManufacturerSpecific",
	CommandClassPowerlevel:                               "Powerlevel",
	CommandClassProtection:                               "Protection",
	CommandClassLock:                                     "Lock",
	CommandClassNodeNaming:                               "NodeNaming",
	CommandClassFirmwareUpdateMd:                         "FirmwareUpdateMd",
	CommandClassGroupingName:                             "GroupingName",
	CommandClassRemoteAssociationActivate:                "RemoteAssociationActivate",
	CommandClassRemoteAssociation:                        "RemoteAssociation",
	CommandClassBattery:                                  "Battery",
	CommandClassClock:                                    "Clock",
	CommandClassHail:                                     "Hail",
	CommandClassWakeUp:                                   "WakeUp",
	CommandClassAssociation:                              "Association",
	CommandClassVersion:                                  "Version",
	CommandClassIndicator:                                "Indicator",
	CommandClassProprietary:                              "Proprietary",
	CommandClassLanguage:                                 "Language",
	CommandClassTime:                                     "Time",
	CommandClassTimeParameters:                           "TimeParameters",
	CommandClassGeographicLocation:                       "GeographicLocation",
	CommandClassMultiChannelAssociationV2:                "MultiChannelAssociationV2",
	CommandClassMultiCmd:                                 "MultiCmd",
	CommandClassEnergyProduction:                         "EnergyProduction",
	CommandClassManufacturerProprietary:                  "ManufacturerProprietary",
	CommandClassScreenMd:                                 "ScreenMd",
	CommandClassScreenAttributes:                         "ScreenAttributes",
	CommandClassSimpleAvControl:                          "SimpleAvControl",
	CommandClassAvContentDirectoryMd:                     "AvContentDirectoryMd",
	CommandClassAvRendererStatus:                         "AvRendererStatus",
	CommandClassAvContentSearchMd:                        "AvContentSearchMd",
	CommandClassSecurity:                                 "Security",
	CommandClassAvTaggingMd:                              "AvTaggingMd",
	CommandClassIpConfiguration:                          "IpConfiguration",
	CommandClassAssociationCommandConfiguration:          "AssociationCommandConfiguration",
	CommandClassSensorAlarm:                              "SensorAlarm",
	CommandClassSilenceAlarm:                             "SilenceAlarm",
	CommandClassSensorConfiguration:                      "SensorConfiguration",
	CommandClassMark:                                     "Mark",
	CommandClassNonInteroperable:                         "NonInteroperable",
}

func GetCommandClassString(cc byte) string {
	if val, ok := CommandClassNames[cc]; ok {
		return val + fmt.Sprintf(" (0x%X)", cc)
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", cc)
	}
}
