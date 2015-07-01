package zwave

const (
	BasicTypeController       = 0x01
	BasicTypeStaticController = 0x02
	BasicTypeSlave            = 0x03
	BasicTypeRoutingSlave     = 0x04
)

var BasicTypeNames map[byte]string = map[byte]string{
	BasicTypeController:       "Controller",
	BasicTypeStaticController: "Static Controller",
	BasicTypeSlave:            "Slave",
	BasicTypeRoutingSlave:     "Routing Slave",
}

func GetBasicTypeName(basicType byte) string {
	if val, ok := BasicTypeNames[basicType]; ok {
		return val
	} else {
		return "Unknown"
	}
}

const (
	GenericTypeGenericController  = 0x01
	GenericTypeStaticController   = 0x02
	GenericTypeAVControlPoint     = 0x03
	GenericTypeDisplay            = 0x04
	GenericTypeNetworkExtender    = 0x05
	GenericTypeAppliance          = 0x06
	GenericTypeSensorNotification = 0x07
	GenericTypeThermostat         = 0x08
	GenericTypeWindowCovering     = 0x09
	GenericTypeSwitchBinary       = 0x10
	GenericTypeSwitchMultilevel   = 0x11
	GenericTypeSwitchRemove       = 0x12
	GenericTypeSwitchToggle       = 0x13
	GenericTypeZipNode            = 0x15
	GenericTypeVentilation        = 0x16
	GenericTypeSecurityPanel      = 0x17
	GenericTypeWallController     = 0x18
	GenericTypeSensorBinary       = 0x20
	GenericTypeSensorMultilevel   = 0x21
	GenericTypeMeterPulse         = 0x30
	GenericTypeMeter              = 0x31
	GenericTypeEntryControl       = 0x40
	GenericTypeSemiInteroperable  = 0x50
	GenericTypeSensorAlarm        = 0xA1
	GenericTypeRepeaterSlave      = 0x0F
	GenericTypeNonInteroperable   = 0xFF
)

var GenericTypeNames map[byte]string = map[byte]string{
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
	GenericTypeSwitchRemove:       "Switch Remove",
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

func GetGenericTypeName(genericType byte) string {
	if val, ok := GenericTypeNames[genericType]; ok {
		return val
	} else {
		return "Unknown"
	}
}
