package ccgen

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

type ValueAttrib struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}

type Variant struct {
	Key         string `xml:"key,attr"`
	ParamOffset byte   `xml:"paramoffs,attr"`
	HasDefines  bool   `xml:"hasdefines,attr"`
	ShowHex     bool   `xml:"showhex,attr"`
	Signed      bool   `xml:"signed,attr"`
	SizeMask    string `xml:"sizemask,attr"`
	SizeOffset  string `xml:"sizeoffs,attr"`
}

type Bitmask struct {
	Key          string `xml:"key,attr"`
	ParamOffset  byte   `xml:"paramoffs,attr"`
	LengthOffset byte   `xml:"lenoffs,attr"`
	LengthMask   byte   `xml:"lenmask,attr"`
	Length       byte   `xml:"len,attr"`
}

type Word struct {
	Key        string `xml:"key,attr"`
	HasDefines bool   `xml:"hasdefines,attr"`
	ShowHex    bool   `xml:"showhex,attr"`
}
