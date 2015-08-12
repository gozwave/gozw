package ccgen

const devicesTemplate = `package commandclass

import "fmt"

type BasicDevice byte

func (b BasicDevice) String() string {
  if val, ok := BasicTypeNames[b]; ok {
		return val + fmt.Sprintf(" (0x%X)", basicType)
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", basicType)
	}
}

type GenericDevice byte

func (g GenericDevice) String() string {
	if val, ok := GenericTypeNames[g]; ok {
		return val + fmt.Sprintf(" (0x%X)", genericType)
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", genericType)
	}
}

type SpecificDevice byte

type DeviceType struct {
  BasicType BasicDevice
  GenericType GenericDevice
  SpecificType SpecificDevice
}

func getSpecificDeviceTypeName(genericType GenericDevice, specificType SpecificDevice) string {
	if val, ok := SpecificTypeNames[genericType][specificType]; ok {
		return val + fmt.Sprintf(" (0x%X)", specificType)
	} else {
		return "Unknown" + fmt.Sprintf(" (0x%X)", specificType)
	}
}

func (d DeviceType) String() string {
  return fmt.Sprintf("%s / %s / %s", d.BasicType, d.GenericType, d.SpecificDeviceString())
}

func (d DeviceType) SpecificDeviceString() string {
  return getSpecificDeviceTypeName(d.GenericType, d.SpecificType)
}

const (
  {{range .BasicDevices}}
  {{ToGoName .Name}} BasicDevice = {{.Key}}
  {{end}}
)

var BasicDeviceNames map[BasicDevice]string = map[BasicDevice]string{
  {{range .BasicDevices}}
  {{ToGoName .Name}}: "{{.Help}}",
  {{end}}
}

const (
  {{range .GenericDevices}}
  {{ToGoName .Name}} GenericDevice = {{.Key}}
  {{end}}
)

var GenericDeviceNames map[GenericDevice]string = map[GenericDevice]string{
  {{range .GenericDevices}}
  {{ToGoName .Name}}: "{{.Help}}",
  {{end}}
}

const (
  SpecificTypeNotUsed SpecificDevice = 0x00
  {{range .GenericDevices}}{{range .SpecificDevices}}{{if NotZeroByte .Key}}
  {{ToGoName .Name}} SpecificDevice = {{.Key}}
  {{end}}{{end}}{{end}}
)

var SpecificTypeNames map[GenericDevice]map[SpecificDevice]string = map[GenericDevice]map[SpecificDevice]string{
  {{range .GenericDevices}}{{ToGoName .Name}}: map[SpecificDevice]string{
    {{range .SpecificDevices}}
    {{ToGoName .Name}}: "{{.Help}}",
    {{end}}
  },
  {{end}}
}
`
