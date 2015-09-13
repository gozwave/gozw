package commandclass

import "encoding"

//go:generate zwgen command-classes -c zwgen.config.yaml -o .
//go:generate zwgen parser -c zwgen.config.yaml -o ./command-classes.gen.go
//go:generate zwgen devices -c zwgen.config.yaml -o ./devices.gen.go

type Command interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	CommandClassID() byte
	CommandID() byte
}
