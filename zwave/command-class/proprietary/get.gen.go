// THIS FILE IS AUTO-GENERATED BY CCGEN
// DO NOT MODIFY

package proprietary

// <no value>

type ProprietaryGet struct {
	Data []byte
}

func ParseProprietaryGet(payload []byte) ProprietaryGet {
	val := ProprietaryGet{}

	i := 2

	val.Data = payload[i:]

	return val
}