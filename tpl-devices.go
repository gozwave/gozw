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
  {{#each BasicDevices}}
  {{toPascalCase Name}} BasicDevice = {{Key}}
  {{/each}}
)

var BasicDeviceNames map[BasicDevice]string = map[BasicDevice]string{
  {{#each BasicDevices}}
  {{toPascalCase Name}}: "{{Help}}",
  {{/each}}
}

const (
  {{#each GenericDevices}}
  {{toPascalCase Name}} GenericDevice = {{Key}}
  {{/each}}
)

var GenericDeviceNames map[GenericDevice]string = map[GenericDevice]string{
  {{#each GenericDevices}}
  {{toPascalCase Name}}: "{{Help}}",
  {{/each}}
}

const (
  SpecificTypeNotUsed SpecificDevice = 0x00
  {{#each GenericDevices}}{{#each SpecificDevices}}{{#notZeroByte Key}}
  {{toPascalCase Name}} SpecificDevice = {{Key}}
  {{/notZeroByte}}{{/each}}{{/each}}
)

var SpecificTypeNames map[GenericDevice]map[SpecificDevice]string = map[GenericDevice]map[SpecificDevice]string{
  {{#each GenericDevices}}{{toPascalCase Name}}: map[SpecificDevice]string{
    {{#each SpecificDevices}}
    {{toPascalCase Name}}: "{{Help}}",
    {{/each}}
  },
  {{/each}}
}
`
