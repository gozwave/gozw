package util

import (
	"encoding/binary"
	"fmt"
	"math"
)

type ZWFloat struct {
	Value float64
	Scale byte
}

func ParseZWFloat(size, scale, precision byte, value []byte) (val *ZWFloat, err error) {
	val = &ZWFloat{}
	val.Scale = scale

	switch size {
	case 1:
		val.Value = float64(int8(value[0])) / math.Pow(10, float64(precision))
		return
	case 2:
		value := int16(binary.BigEndian.Uint16(value))
		val.Value = float64(int16(value)) / math.Pow(10, float64(precision))
		return
	case 4:
		value := int32(binary.BigEndian.Uint32(value))
		val.Value = float64(int32(value)) / math.Pow(10, float64(precision))
		return
	default:
		return nil, fmt.Errorf("Invalid size field: %d", size)
	}
}
