package ccgen

import (
	"encoding/xml"
	"strconv"
	"strings"

	"github.com/reiver/go-stringcase"
)

type ZwClasses struct {
	XMLName        xml.Name        `xml:"zw_classes"`
	BasicDevices   []BasicDevice   `xml:"bas_dev"`
	GenericDevices []GenericDevice `xml:"gen_dev"`
	CommandClasses []CommandClass  `xml:"cmd_class"`
}

type BasicDevice struct {
	Name     string `xml:"name,attr"`
	Key      string `xml:"key,attr"`
	Help     string `xml:"help,attr"`
	ReadOnly bool   `xml:"read_only,attr"`
	Comment  string `xml:"comment,attr"`
}

type GenericDevice struct {
	Name            string           `xml:"name,attr"`
	Key             string           `xml:"key,attr"`
	Help            string           `xml:"help,attr"`
	ReadOnly        bool             `xml:"read_only,attr"`
	Comment         string           `xml:"comment,attr"`
	SpecificDevices []SpecificDevice `xml:"spec_dev"`
}

type SpecificDevice struct {
	Name     string `xml:"name,attr"`
	Key      string `xml:"key,attr"`
	Help     string `xml:"help,attr"`
	ReadOnly bool   `xml:"read_only,attr"`
	Comment  string `xml:"comment,attr"`
}

type CommandClass struct {
	Name     string    `xml:"name,attr"`
	Key      string    `xml:"key,attr"`
	Version  int       `xml:"version,attr"`
	Help     string    `xml:"help,attr"`
	Comment  string    `xml:"comment,attr"`
	Commands []Command `xml:"cmd"`
}

func (c CommandClass) GetDirName() string {
	ccname := strings.Replace(c.Name, "COMMAND_CLASS_", "", 1)
	ccname = stringcase.ToPropertyCase(ccname)

	if c.Version > 1 {
		versionStr := strconv.Itoa(c.Version)
		ccname += "-v" + versionStr
	}

	return ccname
}

func (c CommandClass) GetPackageName() string {
	ccname := strings.Replace(c.Name, "COMMAND_CLASS_", "", 1)
	ccname = stringcase.ToLowerCase(stringcase.ToPascalCase(ccname))

	if c.Version > 1 {
		versionStr := strconv.Itoa(c.Version)
		ccname += "v" + versionStr
	}

	return ccname
}

func (c CommandClass) CanGenerate() (can bool, reason string) {
	if c.Name == "ZWAVE_CMD_CLASS" {
		return false, "Not an actual command class (also stupidly complicated parsing rules)"
	}

	if c.Name == "COMMAND_CLASS_ZIP_6LOWPAN" ||
		c.Name == "COMMAND_CLASS_ZIP_ND" ||
		c.Name == "COMMAND_CLASS_ZIP_GATEWAY" ||
		c.Name == "COMMAND_CLASS_ZIP_PORTAL" ||
		c.Name == "COMMAND_CLASS_ZIP" ||
		c.Name == "COMMAND_CLASS_IP_ASSOCIATION" ||
		c.Name == "COMMAND_CLASS_TRANSPORT_SERVICE" ||
		c.Name == "COMMAND_CLASS_CONTROLLER_REPLICATION" ||
		c.Name == "COMMAND_CLASS_CRC_16_ENCAP" ||
		c.Name == "COMMAND_CLASS_IP_CONFIGURATION" {
		return false, "Skipped (no current intention to support)"
	}

	for _, cmd := range c.Commands {
		for _, param := range cmd.Params {
			if param.Type == "MARKER" {
				return false, "Contains a MARKER"
			}
		}
	}

	return true, ""
}

type Command struct {
	Name     string `xml:"name,attr"`
	Key      string `xml:"key,attr"`
	Type     string `xml:"type,attr"`
	HashCode string
	Comment  string `xml:"comment,attr"`

	Params []Param `xml:"param"`
}
