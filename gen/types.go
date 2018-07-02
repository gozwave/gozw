package gen

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
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

	Enabled bool `xml:"-"`
}

func (c CommandClass) GetBaseName() string {
	return strings.Replace(c.Name, "COMMAND_CLASS_", "", 1)
}

func (c CommandClass) GetConstName() string {
	name := c.GetBaseName()
	if c.Version > 1 {
		versionStr := strconv.Itoa(c.Version)
		name += "_V" + versionStr
	}

	return stringcase.ToPascalCase(name)
}

func (c CommandClass) GetDirName() string {
	ccname := stringcase.ToPropertyCase(c.GetBaseName())

	if c.Version > 1 {
		versionStr := strconv.Itoa(c.Version)
		ccname += "-v" + versionStr
	}

	return ccname
}

func (c CommandClass) GetPackageName() string {
	ccname := stringcase.ToLowerCase(stringcase.ToPascalCase(c.GetBaseName()))

	if c.Version > 1 {
		versionStr := strconv.Itoa(c.Version)
		ccname += "v" + versionStr
	}

	return ccname
}

func (c CommandClass) CanGenerate() (can bool, reason string) {
	if len(c.Commands) == 0 {
		return false, "No commands"
	}

	if c.Name == "ZWAVE_CMD_CLASS" {
		return false, "Not an actual command class"
	}

	if c.Name == "COMMAND_CLASS_ZIP" {
		return false, "Not supported"
	}

	return true, ""
}

func (c CommandClass) CanGen() (can bool) {
	can, _ = c.CanGenerate()
	return
}

type Command struct {
	Name     string `xml:"name,attr"`
	Key      string `xml:"key,attr"`
	Type     string `xml:"type,attr"`
	HashCode string
	Comment  string `xml:"comment,attr"`

	Params []Param `xml:"param"`
}

func (c Command) GetBaseName(cc CommandClass) string {
	commandName := c.Name

	if strings.HasPrefix(strings.ToLower(commandName), "command_") &&
		!strings.HasPrefix(strings.ToLower(commandName), "command_class_") {
		commandName = commandName[8:]
	}

	ccBaseName := cc.GetBaseName()

	if strings.HasPrefix(commandName, ccBaseName) {
		if len(commandName) > len(ccBaseName) {
			commandName = commandName[len(ccBaseName)+1:]
		}
	}

	return commandName
}

func (c Command) GetStructName(cc CommandClass) string {
	return stringcase.ToPascalCase(c.GetBaseName(cc))
}

func (c Command) GetFileName(cc CommandClass) string {
	return stringcase.ToPropertyCase(c.GetBaseName(cc))
}

type Param struct {
	Key            string `xml:"key,attr"`
	Name           string `xml:"name,attr"`
	Type           string `xml:"type,attr"`
	TypeHashCode   string `xml:"typehashcode,attr"`
	Comment        string `xml:"comment,attr"`
	EncapType      string `xml:"encaptype,attr"`
	OptionalOffset string `xml:"optionaloffs,attr"`
	OptionalMask   string `xml:"optionalmask,attr"`
	Encapsulated   bool   `xml:"encapsulated,attr"`
	CommandMask    string `xml:"cmd_mask,attr"`
	Affix          bool   `xml:"affix,attr"`

	ArrayAttrib []ArrayAttrib `xml:"arrayattrib"`
	Bit24       []Bit24       `xml:"bit_24"`
	BitField    []BitField    `xml:"bitfield"`
	BitFlag     []BitFlag     `xml:"bitflag"`
	Const       []Const       `xml:"const"`
	DWord       []DWord       `xml:"dword"`
	FieldEnum   []FieldEnum   `xml:"fieldenum"`
	ValueAttrib []ValueAttrib `xml:"valueattrib"`
	Variant     []Variant     `xml:"variant"`
	Bitmask     []Bitmask     `xml:"bitmask"`
	Word        []Word        `xml:"word"`
}

func (p Param) IsNotReserved() bool {
	return !isReservedString(p.Name)
}

func (p Param) GetEncodedByteLength() (uint8, error) {
	switch p.Type {

	case "BYTE":
		return 1, nil

	case "CONST":
		return 1, nil

	case "MARKER":
		return 1, nil

	case "STRUCT_BYTE":
		return 1, errors.New("Unimplemented param type: STRUCT_BYTE")

	case "WORD":
		return 2, nil

	case "BIT_24":
		return 3, nil

	case "DWORD":
		return 4, nil

	case "ARRAY":
		if len(p.ArrayAttrib) > 0 && p.ArrayAttrib[0].Length != 0 {
			return byte(p.ArrayAttrib[0].Length), nil
		} else {
			return 0, errors.New("Field has unknown or indeterminate length")
		}

	default:
		fmt.Println(p.Name, p.Type)
		return 0, errors.New("Field has unknown or indeterminate length")

	}
}

func (p Param) GetGoType() (string, error) {
	switch p.Type {

	case "ARRAY":
		if p.ArrayAttrib != nil && len(p.ArrayAttrib) == 1 {
			if p.ArrayAttrib[0].IsAscii {
				return "string", nil
			} else {
				return "[]byte", nil
			}
		} else if p.ArrayAttrib != nil {
			return "", errors.New("Weird number of <ArrayAttrib> elements")
		}
		return "[]byte", nil

	case "BIT_24":
		return "uint32", nil

	case "BITMASK":
		// @todo there are some command classes (that we currently don't generate)
		// that have a BITMASK in the middle of the payload. We don't currently
		// support that.
		return "[]byte", nil

	case "BYTE":
		return "byte", nil

	case "CONST":
		return "byte", nil

	case "DWORD":
		return "uint32", nil

	case "ENUM":
		return "", errors.New("Unimplemented param type: ENUM")

	case "ENUM_ARRAY":
		return "", errors.New("Unimplemented param type: ENUM_ARRAY")

	case "MARKER":
		return "", errors.New("Unimplemented param type: MARKER")

	case "MULTI_ARRAY":
		return "", errors.New("Unimplemented param type: MULTI_ARRAY")

	case "STRUCT_BYTE":
		return "", errors.New("Unimplemented param type: STRUCT_BYTE")

	case "VARIANT":
		return "[]byte", nil

	case "WORD":
		return "uint16", nil

	default:
		return "", errors.New("Unknown param type: " + p.Type)

	}
}

type ArrayAttrib struct {
	Key     string `xml:"key,attr"`
	Length  int    `xml:"len,attr"`
	IsAscii bool   `xml:"is_ascii,attr"`
}

type ArrayLen struct {
	Key          string `xml:"key,attr"`
	ParamOffset  byte   `xml:"paramoffs,attr"`
	LengthOffset byte   `xml:"lenoffs,attr"`
	LengthMask   byte   `xml:"lenmask,attr"`
}

type Bit24 struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}

type BitField struct {
	Key       string `xml:"key,attr"`
	FieldName string `xml:"fieldname,attr"`
	FieldMask string `xml:"fieldmask,attr"`
	Shifter   uint8  `xml:"shifter,attr"`
}

func (b BitField) IsNotReserved() bool {
	return !isReservedString(b.FieldName)
}

type BitFlag struct {
	Key      string `xml:"key,attr"`
	FlagName string `xml:"flagname,attr"`
	FlagMask string `xml:"flagmask,attr"`
}

func (b BitFlag) IsNotReserved() bool {
	return !isReservedString(b.FlagName)
}

type Const struct {
	Key      string `xml:"key,attr"`
	FlagName string `xml:"flagname,attr"`
	FlagMask string `xml:"flagmask,attr"`
}

type DWord struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}

type Enum struct {
	Key  string `xml:"key,attr"`
	Name string `xml:"name,attr"`
}

type FieldEnum struct {
	Key       string `xml:"key,attr"`
	FieldName string `xml:"fieldname,attr"`
	FieldMask string `xml:"fieldmask,attr"`
	Shifter   uint8  `xml:"shifter,attr"`
	Value     string `xml:"value,attr"`

	EnumValues []FieldEnum `xml:"fieldenum"`
}

func (f FieldEnum) IsNotReserved() bool {
	return !isReservedString(f.FieldName)
}

type ValueAttrib struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}

type Variant struct {
	Key         string `xml:"key,attr"`
	ParamOffset Byte   `xml:"paramoffs,attr"`
	HasDefines  bool   `xml:"hasdefines,attr"`
	ShowHex     bool   `xml:"showhex,attr"`
	Signed      bool   `xml:"signed,attr"`
	SizeMask    string `xml:"sizemask,attr"`
	SizeOffset  string `xml:"sizeoffs,attr"`

	MarkerDelimited bool
	MarkerValue     string
	RemainingBytes  uint8
}

type Bitmask struct {
	Key          string `xml:"key,attr"`
	ParamOffset  Byte   `xml:"paramoffs,attr"`
	LengthOffset Byte   `xml:"lenoffs,attr"`
	LengthMask   Byte   `xml:"lenmask,attr"`
	Length       Byte   `xml:"len,attr"`
}

type Word struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}

type Byte byte

func (b *Byte) UnmarshalXMLAttr(attr xml.Attr) error {
	if strings.HasPrefix(attr.Value, "0x") {
		bs, err := hex.DecodeString(attr.Value[2:])
		if err != nil {
			return err
		}
		*b = Byte(bs[0])
	} else {
		n, err := strconv.ParseUint(attr.Value, 10, 8)
		if err != nil {
			return err
		}
		*b = Byte(n)
	}
	return nil
}
