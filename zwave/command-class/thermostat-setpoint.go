package commandclass

import (
	"encoding/binary"
	"errors"
	"math"
)

/* Thermostat Setpoint command class commands */
const (
	CommandThermostatSetpointSet             byte = 0x01
	CommandThermostatSetpointGet                  = 0x02
	CommandThermostatSetpointReport               = 0x03
	CommandThermostatSetpointSupportedGet         = 0x04
	CommandThermostatSetpointSupportedReport      = 0x05
)

const ThermostatSetpointTypeMask byte = 0x0F

type ThermostatSetpointType byte

const (
	ThermostatSetpointTypeNotSupported   ThermostatSetpointType = 0x00
	ThermostatSetpointTypeHeating                               = 0x01
	ThermostatSetpointTypeCooling                               = 0x02
	ThermostatSetpointTypeNotSupported1                         = 0x03
	ThermostatSetpointTypeNotSupported2                         = 0x04
	ThermostatSetpointTypeNotSupported3                         = 0x05
	ThermostatSetpointTypeNotSupported4                         = 0x06
	ThermostatSetpointTypeFurnace                               = 0x07
	ThermostatSetpointTypeDryAir                                = 0x08
	ThermostatSetpointTypeMoistAir                              = 0x09
	ThermostatSetpointTypeAutoChangeover                        = 0x0A
)

type ThermostatSetpointScale byte

const (
	SetpointScaleCelcius   byte = 0x0
	SetpointScaleFarenheit      = 0x1
)

const (
	setpointTypeMask       byte = 0x0F
	setpointPrecisionMask       = 0xE0
	setpointPrecisionShift      = 0x05
	setpointScaleMask           = 0x18
	setpointScaleShift          = 0x03
	setpointSizeMask            = 0x07
)

type Temperature struct {
	Value float64
	Scale ThermostatSetpointScale
}

type ThermostatSetpoint struct {
	Type      ThermostatSetpointType
	Precision byte
	Scale     byte
	Size      byte
	Value     []byte
}

// Even though we can handle receiving any theoretically-valid value from a
// thermostat, we're only going to support 0-100 degrees (C or F) for setting,
// and for now, we're not going to support decimal values.
func NewThermostatSetpoint(setpointType ThermostatSetpointType, temp Temperature) (*ThermostatSetpoint, error) {
	val := math.Floor(temp.Value)
	if val > 100 {
		return nil, errors.New("Setpoint temperature too high")
	}

	if val < 0 {
		return nil, errors.New("Setpoint temperature too low")
	}

	setpoint := &ThermostatSetpoint{
		Type:      setpointType,
		Precision: 0,
		Scale:     byte(temp.Scale),
		Size:      1,
		Value:     []byte{byte(val)},
	}

	return setpoint, nil
}

func (t *ThermostatSetpoint) GetTemperature() (*Temperature, error) {
	switch t.Size {
	case 1:
		return &Temperature{
			Value: float64(int8(t.Value[0])) / math.Pow(10, float64(t.Precision)),
			Scale: ThermostatSetpointScale(t.Scale),
		}, nil
	case 2:
		value := int16(binary.BigEndian.Uint16(t.Value))
		return &Temperature{
			Value: float64(int16(value)) / math.Pow(10, float64(t.Precision)),
			Scale: ThermostatSetpointScale(t.Scale),
		}, nil
	case 4:
		value := int32(binary.BigEndian.Uint32(t.Value))
		return &Temperature{
			Value: float64(int32(value)) / math.Pow(10, float64(t.Precision)),
			Scale: ThermostatSetpointScale(t.Scale),
		}, nil
	default:
		return nil, errors.New("Invalid size field in setpoint report")
	}
}

func NewThermostatSetpointSet(setpointType ThermostatSetpointType, temperature Temperature) ([]byte, error) {
	setpoint, err := NewThermostatSetpoint(setpointType, temperature)
	if err != nil {
		return nil, err
	}

	precision := (setpoint.Precision << setpointPrecisionShift) & setpointPrecisionMask
	scale := (setpoint.Scale << setpointScaleShift) & setpointScaleMask
	size := setpoint.Size & setpointSizeMask

	payload := []byte{
		CommandClassThermostatSetpoint,
		CommandThermostatSetpointSet,
		byte(setpointType) & setpointTypeMask, // setpoint type
		precision | scale | size,
	}

	payload = append(payload, setpoint.Value...)

	return payload, nil
}

func ParseThermostatSetpointReport(payload []byte) ThermostatSetpoint {
	return ThermostatSetpoint{
		Type:      ThermostatSetpointType(payload[2] & 0x0F),
		Precision: (payload[3] & setpointPrecisionMask) >> setpointPrecisionShift,
		Scale:     (payload[3] & setpointScaleMask) >> setpointScaleShift,
		Size:      (payload[3] & setpointSizeMask),
		Value:     payload[4:],
	}
}
