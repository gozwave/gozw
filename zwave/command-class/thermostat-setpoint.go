package commandclass

import (
	"encoding/binary"
	"errors"
	"math"
)

/* Thermostat Setpoint command class commands */
const (
	CommandThermostatSetpointVersion         byte = 0x01
	CommandThermostatSetpointGet                  = 0x02
	CommandThermostatSetpointReport               = 0x03
	CommandThermostatSetpointSet                  = 0x01
	CommandThermostatSetpointSupportedGet         = 0x04
	CommandThermostatSetpointSupportedReport      = 0x05
)

const (
	ThermostatSetpointTypeMask           byte = 0x0F
	ThermostatSetpointTypeNotSupported        = 0x00
	ThermostatSetpointTypeHeating             = 0x01
	ThermostatSetpointTypeCooling             = 0x02
	ThermostatSetpointTypeNotSupported1       = 0x03
	ThermostatSetpointTypeNotSupported2       = 0x04
	ThermostatSetpointTypeNotSupported3       = 0x05
	ThermostatSetpointTypeNotSupported4       = 0x06
	ThermostatSetpointTypeFurnace             = 0x07
	ThermostatSetpointTypeDryAir              = 0x08
	ThermostatSetpointTypeMoistAir            = 0x09
	ThermostatSetpointTypeAutoChangeover      = 0x0A
)

const (
	SetpointScaleCelcius   byte = 0x0
	SetpointScaleFarenheit      = 0x1
)

const (
	setpointPrecisionMask  byte = 0xE0
	setpointPrecisionShift      = 0x05
)

const (
	thermostatSetpointTypeMask       byte = 0x0F
	thermostatSetpointPrecisionMask       = 0xE0
	thermostatSetpointPrecisionShift      = 0x05
	thermostatSetpointScaleMask           = 0x18
	thermostatSetpointScaleShift          = 0x03
	thermostatSetpointSizeMask            = 0x07
)

type Temperature struct {
	Value float64
	Scale byte
}

type ThermostatSetpointReport struct {
	Type      byte
	Precision byte
	Scale     byte
	Size      byte
	Value     []byte
}

func ParseThermostatSetpointReport(payload []byte) ThermostatSetpointReport {
	return ThermostatSetpointReport{
		Type:      payload[2] & 0x0F,
		Precision: (payload[3] & 0xE0) >> 5,
		Scale:     (payload[3] & 0x18) >> 3,
		Size:      (payload[3] & 0x07),
		Value:     payload[4:],
	}
}

func (t *ThermostatSetpointReport) GetTemperature() (*Temperature, error) {
	switch t.Size {
	case 1:
		return &Temperature{
			Value: float64(int8(t.Value[0])) / math.Pow(10, float64(t.Precision)),
			Scale: t.Scale,
		}, nil
	case 2:
		value := int16(binary.BigEndian.Uint16(t.Value))
		return &Temperature{
			Value: float64(int16(value)) / math.Pow(10, float64(t.Precision)),
			Scale: t.Scale,
		}, nil
	case 4:
		value := int32(binary.BigEndian.Uint32(t.Value))
		return &Temperature{
			Value: float64(int32(value)) / math.Pow(10, float64(t.Precision)),
			Scale: t.Scale,
		}, nil
	default:
		return nil, errors.New("Invalid size field in setpoint report")
	}
}
