package cc

import "fmt"

type BasicDeviceType byte

func (b BasicDeviceType) String() string {
	if val, ok := BasicDeviceTypeNames[b]; ok {
		return val + fmt.Sprintf(" (0x%X)", byte(b))
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", byte(b))
	}
}

type GenericDeviceType byte

func (g GenericDeviceType) String() string {
	if val, ok := GenericDeviceTypeNames[g]; ok {
		return val + fmt.Sprintf(" (0x%X)", byte(g))
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", byte(g))
	}
}

type SpecificDeviceType byte

type DeviceType struct {
	BasicType    BasicDeviceType
	GenericType  GenericDeviceType
	SpecificType SpecificDeviceType
}

func getSpecificDeviceTypeName(genericType GenericDeviceType, specificType SpecificDeviceType) string {
	if val, ok := SpecificDeviceTypeNames[genericType][specificType]; ok {
		return val + fmt.Sprintf(" (0x%X)", byte(specificType))
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", byte(specificType))
	}
}

func (d DeviceType) String() string {
	return fmt.Sprintf("%s / %s / %s", d.BasicType, d.GenericType, d.SpecificDeviceTypeString())
}

func (d DeviceType) SpecificDeviceTypeString() string {
	return getSpecificDeviceTypeName(d.GenericType, d.SpecificType)
}
