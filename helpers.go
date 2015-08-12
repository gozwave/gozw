package ccgen

import (
	"strings"

	"github.com/aymerick/raymond"
	"github.com/reiver/go-stringcase"
)

func getGoType(paramType string) string {
	switch paramType {
	case "ARRAY":
		return "[]byte"
	case "BIT_24":
		panic("Unimplemented param type: BIT_24")
	case "BITMASK":
		return "byte"
	case "BYTE":
		return "byte"
	case "CONST":
		return "byte"
	case "DWORD":
		return "uint32"
	case "ENUM":
		return "byte"
	case "ENUM_ARRAY":
		return "[]byte"
	case "MARKER":
		panic("Unimplemented param type: MARKER")
	case "MULTI_ARRAY":
		panic("Unimplemented param type: MULTI_ARRAY")
	case "STRUCT_BYTE":
		panic("Unimplemented param type: STRUCT_BYTE")
	case "VARIANT":
		return "[]byte"
	case "WORD":
		return "uint16"
	default:
		panic("Unknown param type: " + paramType)
	}
}

// func notZeroByte(opts *raymond.Options) bool {
// 	str := opts.ParamStr(0)
// 	return str == "0x00"
// }

func notZeroByte(str string, options *raymond.Options) string {
	if str != "0x00" {
		return options.Fn()
	}
	return ""
}

func toPackageName(ccname string) string {
	ccname = strings.Replace(ccname, "COMMAND_CLASS_", "", 1)
	return stringcase.ToLowerCase(stringcase.ToPascalCase(ccname))
}
