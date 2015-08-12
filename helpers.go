package ccgen

import (
	"regexp"

	"github.com/reiver/go-stringcase"
)

var nameReplacer = regexp.MustCompile(`[^\w]+`)

func toGoName(name string) string {
	return nameReplacer.ReplaceAllString(stringcase.ToPascalCase(name), "")
}

// func notZeroByte(opts *raymond.Options) bool {
// 	str := opts.ParamStr(0)
// 	return str == "0x00"
// }

func notZeroByte(str string) string {
	if str != "0x00" {
		return str
	}
	return ""
}
