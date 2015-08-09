package commandclass

import "encoding/binary"

const (
	CommandManufacturerSpecificGet    byte = 0x04
	CommandManufacturerSpecificReport      = 0x05
)

type ManufacturerSpecificReport struct {
	ManufacturerID uint16
	ProductTypeID  uint16
	ProductID      uint16
}

func ParseManufacturerSpecificReport(payload []byte) ManufacturerSpecificReport {
	return ManufacturerSpecificReport{
		ManufacturerID: binary.BigEndian.Uint16(payload[2:4]),
		ProductTypeID:  binary.BigEndian.Uint16(payload[4:6]),
		ProductID:      binary.BigEndian.Uint16(payload[6:8]),
	}
}
