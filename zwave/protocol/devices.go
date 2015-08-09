package protocol

import "fmt"

const (
	BasicTypeController       byte = 0x01
	BasicTypeStaticController      = 0x02
	BasicTypeSlave                 = 0x03
	BasicTypeRoutingSlave          = 0x04
)

var BasicTypeNames = map[byte]string{
	BasicTypeController:       "Controller",
	BasicTypeStaticController: "Static Controller",
	BasicTypeSlave:            "Slave",
	BasicTypeRoutingSlave:     "Routing Slave",
}

const (
	GenericTypeGenericController  byte = 0x01
	GenericTypeStaticController        = 0x02
	GenericTypeAVControlPoint          = 0x03
	GenericTypeDisplay                 = 0x04
	GenericTypeNetworkExtender         = 0x05
	GenericTypeAppliance               = 0x06
	GenericTypeSensorNotification      = 0x07
	GenericTypeThermostat              = 0x08
	GenericTypeWindowCovering          = 0x09
	GenericTypeSwitchBinary            = 0x10
	GenericTypeSwitchMultilevel        = 0x11
	GenericTypeSwitchRemote            = 0x12
	GenericTypeSwitchToggle            = 0x13
	GenericTypeZipNode                 = 0x15
	GenericTypeVentilation             = 0x16
	GenericTypeSecurityPanel           = 0x17
	GenericTypeWallController          = 0x18
	GenericTypeSensorBinary            = 0x20
	GenericTypeSensorMultilevel        = 0x21
	GenericTypeMeterPulse              = 0x30
	GenericTypeMeter                   = 0x31
	GenericTypeEntryControl            = 0x40
	GenericTypeSemiInteroperable       = 0x50
	GenericTypeSensorAlarm             = 0xA1
	GenericTypeRepeaterSlave           = 0x0F
	GenericTypeNonInteroperable        = 0xFF
)

var GenericTypeNames = map[byte]string{
	GenericTypeGenericController:  "Generic Controller",
	GenericTypeStaticController:   "Static Controller",
	GenericTypeAVControlPoint:     "AV Control Point",
	GenericTypeDisplay:            "Display",
	GenericTypeNetworkExtender:    "Network Extender",
	GenericTypeAppliance:          "Appliance",
	GenericTypeSensorNotification: "Sensor Notification",
	GenericTypeThermostat:         "Thermostat",
	GenericTypeWindowCovering:     "Window Covering",
	GenericTypeSwitchBinary:       "Switch Binary",
	GenericTypeSwitchMultilevel:   "Switch Multi-level",
	GenericTypeSwitchRemote:       "Switch Remote",
	GenericTypeSwitchToggle:       "Switch Toggle",
	GenericTypeZipNode:            "Z/IP Node",
	GenericTypeVentilation:        "Ventilation",
	GenericTypeSecurityPanel:      "Security Panel",
	GenericTypeWallController:     "Wall Controller",
	GenericTypeSensorBinary:       "Sensor Binary",
	GenericTypeSensorMultilevel:   "Sensor Multi-level",
	GenericTypeMeterPulse:         "Pulse Meter",
	GenericTypeMeter:              "Meter",
	GenericTypeEntryControl:       "Entry Control",
	GenericTypeSemiInteroperable:  "Semi-Interoperable",
	GenericTypeSensorAlarm:        "Sensor Alarm",
	GenericTypeRepeaterSlave:      "Repeater Slave",
	GenericTypeNonInteroperable:   "Non-Interoperable",
}

const (
	SpecificTypeNotUsed byte = 0x00

	SpecificTypeDoorbell            = 0x12
	SpecificTypeSatelliteReceiver   = 0x04
	SpecificTypeSatelliteReceiverV2 = 0x11

	SpecificTypeSimpleDisplay = 0x01

	SpecificTypeDoorLock                     = 0x01
	SpecificTypeAdvancedDoorLock             = 0x02
	SpecificTypeSecureKeypadDoorLock         = 0x03
	SpecificTypeSecureKeypadDoorLockDeadbolt = 0x04
	SpecificTypeSecureDoor                   = 0x05
	SpecificTypeSecureGate                   = 0x06
	SpecificTypeSecureBarrierAddon           = 0x07
	SpecificTypeSecureBarrierOpenOnly        = 0x08
	SpecificTypeSecureBarrierCloseOnly       = 0x09
	SpecificTypeSecureLockbox                = 0x0A

	SpecificTypePortableRemoteController = 0x01
	SpecificTypePortableSceneController  = 0x02
	SpecificTypePortableInstallerTool    = 0x03
	SpecificTypeRemoteControlAv          = 0x04
	SpecificTypeRemoteControlSimple      = 0x06

	SpecificTypeSimpleMeter          = 0x01
	SpecificTypeAdvEnergyControl     = 0x02
	SpecificTypeWholeHomeMeterSimple = 0x03

	SpecificTypeRepeaterSlave = 0x01

	SpecificTypeZonedSecurityPanel = 0x01

	SpecificTypeEnergyProduction = 0x01

	SpecificTypeAdvZensorNetAlarmSensor   = 0x05
	SpecificTypeAdvZensorNetSmokeSensor   = 0x0A
	SpecificTypeBasicRoutingAlarmSensor   = 0x01
	SpecificTypeBasicRoutingSmokeSensor   = 0x06
	SpecificTypeBasicZensorNetAlarmSensor = 0x03
	SpecificTypeBasicZensorNetSmokeSensor = 0x08
	SpecificTypeRoutingAlarmSensor        = 0x02
	SpecificTypeRoutingSmokeSensor        = 0x07
	SpecificTypeZensorNetAlarmSensor      = 0x04
	SpecificTypeZensorNetSmokeSensor      = 0x09
	SpecificTypeAlarmSensor               = 0x0B

	SpecificTypeRoutingSensorBinary = 0x01

	SpecificTypeRoutingSensorMultilevel = 0x01
	SpecificTypeChimneyFan              = 0x02

	SpecificTypePcController        = 0x01
	SpecificTypeSceneController     = 0x02
	SpecificTypeStaticInstallerTool = 0x03
	SpecificTypeSetTopBox           = 0x04
	SpecificTypeSubSystemController = 0x05
	SpecificTypeTv                  = 0x06
	SpecificTypeGateway             = 0x07

	SpecificTypePowerSwitchBinary  = 0x01
	SpecificTypeSceneSwitchBinary  = 0x03
	SpecificTypePowerStrip         = 0x04
	SpecificTypeSiren              = 0x05
	SpecificTypeValveOpenClose     = 0x06
	SpecificTypeColorTunableBinary = 0x02

	SpecificTypeClassAMotorControl     = 0x05
	SpecificTypeClassBMotorControl     = 0x06
	SpecificTypeClassCMotorControl     = 0x07
	SpecificTypeMotorMultiposition     = 0x03
	SpecificTypePowerSwitchMultilevel  = 0x01
	SpecificTypeSceneSwitchMultilevel  = 0x04
	SpecificTypeFanSwitch              = 0x08
	SpecificTypeColorTunableMultilevel = 0x02

	SpecificTypeSwitchRemoteBinary           = 0x01
	SpecificTypeSwitchRemoteMultilevel       = 0x02
	SpecificTypeSwitchRemoteToggleBinary     = 0x03
	SpecificTypeSwitchRemoteToggleMultilevel = 0x04

	SpecificTypeSwitchToggleBinary     = 0x01
	SpecificTypeSwitchToggleMultilevel = 0x02

	SpecificTypeSetbackScheduleThermostat = 0x03
	SpecificTypeSetbackThermostat         = 0x05
	SpecificTypeSetpointThermostat        = 0x04
	SpecificTypeThermostatGeneral         = 0x02
	SpecificTypeThermostatGeneralV2       = 0x06
	SpecificTypeThermostatHeating         = 0x01

	SpecificTypeResidentialHrv = 0x01

	SpecificTypeSimpleWindowCovering = 0x01

	SpecificTypeZipAdvNode = 0x02
	SpecificTypeZipTunNode = 0x01

	SpecificTypeBasicWallController = 0x01

	SpecificTypeSecureExtender = 0x01

	SpecificTypeGeneralAppliance = 0x01
	SpecificTypeKitchenAppliance = 0x02
	SpecificTypeLaundryAppliance = 0x03

	SpecificTypeNotificationSensor = 0x01
)

var SpecificTypeNames = map[byte]map[byte]string{
	GenericTypeAVControlPoint: map[byte]string{
		SpecificTypeNotUsed:             "Unspecified",
		SpecificTypeDoorbell:            "Doorbell",
		SpecificTypeSatelliteReceiver:   "Satellite Receiver",
		SpecificTypeSatelliteReceiverV2: "Satellite Receiver V2",
	},
	GenericTypeDisplay: map[byte]string{
		SpecificTypeNotUsed:       "Unspecified",
		SpecificTypeSimpleDisplay: "Simple Display",
	},
	GenericTypeEntryControl: map[byte]string{
		SpecificTypeNotUsed:                      "Unspecified",
		SpecificTypeDoorLock:                     "Door Lock",
		SpecificTypeAdvancedDoorLock:             "Advanced Door Lock",
		SpecificTypeSecureKeypadDoorLock:         "Secure Keypad Door Lock",
		SpecificTypeSecureKeypadDoorLockDeadbolt: "Secure Keypad Door Lock Deadbolt",
		SpecificTypeSecureDoor:                   "Secure Door",
		SpecificTypeSecureGate:                   "Secure Gate",
		SpecificTypeSecureBarrierAddon:           "Secure Barrier Addon",
		SpecificTypeSecureBarrierOpenOnly:        "Secure Barrier Open Only",
		SpecificTypeSecureBarrierCloseOnly:       "Secure Barrier Close Only",
		SpecificTypeSecureLockbox:                "Secure Lockbox",
	},
	GenericTypeGenericController: map[byte]string{
		SpecificTypeNotUsed:                  "Unspecified",
		SpecificTypePortableRemoteController: "Portable Remote Controller",
		SpecificTypePortableSceneController:  "Portable Scene Controller",
		SpecificTypePortableInstallerTool:    "Portable Installer Tool",
		SpecificTypeRemoteControlAv:          "Remote Control AV",
		SpecificTypeRemoteControlSimple:      "Remote Control Simple",
	},
	GenericTypeMeter: map[byte]string{
		SpecificTypeNotUsed:              "Unspecified",
		SpecificTypeSimpleMeter:          "Simple Meter",
		SpecificTypeAdvEnergyControl:     "Advanced Energy Control",
		SpecificTypeWholeHomeMeterSimple: "Whole Home Meter Simple",
	},
	GenericTypeMeterPulse: map[byte]string{
		SpecificTypeNotUsed: "Unspecified",
	},
	GenericTypeNonInteroperable: map[byte]string{
		SpecificTypeNotUsed: "Unspecified",
	},
	GenericTypeRepeaterSlave: map[byte]string{
		SpecificTypeNotUsed:       "Unspecified",
		SpecificTypeRepeaterSlave: "Repeater Slave",
	},
	GenericTypeSecurityPanel: map[byte]string{
		SpecificTypeNotUsed:            "Unspecified",
		SpecificTypeZonedSecurityPanel: "Zoned Security Panel",
	},
	GenericTypeSemiInteroperable: map[byte]string{
		SpecificTypeNotUsed:          "Unspecified",
		SpecificTypeEnergyProduction: "Energy Production",
	},
	GenericTypeSensorAlarm: map[byte]string{
		SpecificTypeNotUsed:                   "Unspecified",
		SpecificTypeAdvZensorNetAlarmSensor:   "Adv Zensor Net Alarm Sensor",
		SpecificTypeAdvZensorNetSmokeSensor:   "Adv Zensor Net Smoke Sensor",
		SpecificTypeBasicRoutingAlarmSensor:   "Basic Routing Alarm Sensor",
		SpecificTypeBasicRoutingSmokeSensor:   "Basic Routing Smoke Sensor",
		SpecificTypeBasicZensorNetAlarmSensor: "Basic Zensor Net Alarm Sensor",
		SpecificTypeBasicZensorNetSmokeSensor: "Basic Zensor Net Smoke Sensor",
		SpecificTypeRoutingAlarmSensor:        "Routing Alarm Sensor",
		SpecificTypeRoutingSmokeSensor:        "Routing Smoke Sensor",
		SpecificTypeZensorNetAlarmSensor:      "Zensor Net Alarm Sensor",
		SpecificTypeZensorNetSmokeSensor:      "Zensor Net Smoke Sensor",
		SpecificTypeAlarmSensor:               "Alarm Sensor",
	},
	GenericTypeSensorBinary: map[byte]string{
		SpecificTypeNotUsed:             "Unspecified",
		SpecificTypeRoutingSensorBinary: "Routing Sensor Binary",
	},
	GenericTypeSensorMultilevel: map[byte]string{
		SpecificTypeNotUsed:                 "Unspecified",
		SpecificTypeRoutingSensorMultilevel: "Routing Sensor Multilevel",
		SpecificTypeChimneyFan:              "Chimney Fan",
	},
	GenericTypeStaticController: map[byte]string{
		SpecificTypeNotUsed:             "Unspecified",
		SpecificTypePcController:        "PC Controller",
		SpecificTypeSceneController:     "Scene Controller",
		SpecificTypeStaticInstallerTool: "Static Installer Tool",
		SpecificTypeSetTopBox:           "Set-Top Box",
		SpecificTypeSubSystemController: "Sub-System Controller",
		SpecificTypeTv:                  "TV",
		SpecificTypeGateway:             "Gateway",
	},
	GenericTypeSwitchBinary: map[byte]string{
		SpecificTypeNotUsed:            "Unspecified",
		SpecificTypePowerSwitchBinary:  "Power Switch Binary",
		SpecificTypeSceneSwitchBinary:  "Scene Switch Binary",
		SpecificTypePowerStrip:         "Power Strip",
		SpecificTypeSiren:              "Siren",
		SpecificTypeValveOpenClose:     "Valve Open Close",
		SpecificTypeColorTunableBinary: "Color Tunable Binary",
	},
	GenericTypeSwitchMultilevel: map[byte]string{
		SpecificTypeNotUsed:                "Unspecified",
		SpecificTypeClassAMotorControl:     "Class A Motor Control",
		SpecificTypeClassBMotorControl:     "Class B Motor Control",
		SpecificTypeClassCMotorControl:     "Class C Motor Control",
		SpecificTypeMotorMultiposition:     "Motor Multiposition",
		SpecificTypePowerSwitchMultilevel:  "Power Switch Multilevel",
		SpecificTypeSceneSwitchMultilevel:  "Scene Switch Multilevel",
		SpecificTypeFanSwitch:              "Fan Switch",
		SpecificTypeColorTunableMultilevel: "Color Tunable Multilevel",
	},
	GenericTypeSwitchRemote: map[byte]string{
		SpecificTypeNotUsed:                      "Unspecified",
		SpecificTypeSwitchRemoteBinary:           "Switch Remote Binary",
		SpecificTypeSwitchRemoteMultilevel:       "Switch Remote Multilevel",
		SpecificTypeSwitchRemoteToggleBinary:     "Switch Remote Toggle Binary",
		SpecificTypeSwitchRemoteToggleMultilevel: "Switch Remote Toggle Multilevel",
	},
	GenericTypeSwitchToggle: map[byte]string{
		SpecificTypeNotUsed:                "Unspecified",
		SpecificTypeSwitchToggleBinary:     "Switch Toggle Binary",
		SpecificTypeSwitchToggleMultilevel: "Switch Toggle Multilevel",
	},
	GenericTypeThermostat: map[byte]string{
		SpecificTypeNotUsed:                   "Unspecified",
		SpecificTypeSetbackScheduleThermostat: "Setback Schedule Thermostat",
		SpecificTypeSetbackThermostat:         "Setback Thermostat",
		SpecificTypeSetpointThermostat:        "Setpoint Thermostat",
		SpecificTypeThermostatGeneral:         "Thermostat General",
		SpecificTypeThermostatGeneralV2:       "Thermostat General V2",
		SpecificTypeThermostatHeating:         "Thermostat Heating",
	},
	GenericTypeVentilation: map[byte]string{
		SpecificTypeNotUsed:        "Unspecified",
		SpecificTypeResidentialHrv: "Residential HRV",
	},
	GenericTypeWindowCovering: map[byte]string{
		SpecificTypeNotUsed:              "Unspecified",
		SpecificTypeSimpleWindowCovering: "Simple Window Covering",
	},
	GenericTypeZipNode: map[byte]string{
		SpecificTypeNotUsed:    "Unspecified",
		SpecificTypeZipAdvNode: "Zip Adv Node",
		SpecificTypeZipTunNode: "Zip Tun Node",
	},
	GenericTypeWallController: map[byte]string{
		SpecificTypeNotUsed:             "Unspecified",
		SpecificTypeBasicWallController: "Basic Wall Controller",
	},
	GenericTypeNetworkExtender: map[byte]string{
		SpecificTypeNotUsed:        "Unspecified",
		SpecificTypeSecureExtender: "Secure Extender",
	},
	GenericTypeAppliance: map[byte]string{
		SpecificTypeNotUsed:          "Unspecified",
		SpecificTypeGeneralAppliance: "General Appliance",
		SpecificTypeKitchenAppliance: "Kitchen Appliance",
		SpecificTypeLaundryAppliance: "Laundry Appliance",
	},
	GenericTypeSensorNotification: map[byte]string{
		SpecificTypeNotUsed:            "Unspecified",
		SpecificTypeNotificationSensor: "Notification Sensor",
	},
}

func GetBasicDeviceTypeName(basicType byte) string {
	if val, ok := BasicTypeNames[basicType]; ok {
		return val + fmt.Sprintf(" (0x%X)", basicType)
	}

	return "Unknown" + fmt.Sprintf(" (0x%X)", basicType)
}

func GetGenericDeviceTypeName(genericType byte) string {
	if val, ok := GenericTypeNames[genericType]; ok {
		return val + fmt.Sprintf(" (0x%X)", genericType)
	}

	return "Unknown" + fmt.Sprintf(" (0x%X)", genericType)
}

func GetSpecificDeviceTypeName(genericType byte, specificType byte) string {
	if val, ok := SpecificTypeNames[genericType][specificType]; ok {
		return val + fmt.Sprintf(" (0x%X)", specificType)
	}

	return "Unknown" + fmt.Sprintf(" (0x%X)", specificType)
}
