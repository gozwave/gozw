package commandclass

import "encoding"

//go:generate zwgen -o .

type Command interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

	CommandClassID() byte
	CommandID() byte
}
