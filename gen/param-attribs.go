package gen

import (
	"encoding/hex"
	"encoding/xml"
	"strconv"
	"strings"
)


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
