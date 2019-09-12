package main

import (
	"regexp"
	"strings"

	stringcase "github.com/reiver/go-stringcase"
)

var nameReplacer = regexp.MustCompile(`[^\w]+`)

func isLastKey(key string, keys []string) bool {
	return keys[len(keys)-1] == key
}

func toGoName(name string) string {
	return nameReplacer.ReplaceAllString(stringcase.ToPascalCase(name), "")
}

func toGoNameLower(name string) string {
	return nameReplacer.ReplaceAllString(stringcase.ToCamelCase(name), "")
}

func notZeroByte(str string) string {
	if str != "0x00" {
		return str
	}
	return ""
}

func isReservedString(str string) bool {
	str = strings.ToLower(str)
	if strings.HasPrefix(str, "reserved") {
		return true
	}

	if str == "res" {
		return true
	}

	return false
}
