package cc

import (
	"encoding"
	"errors"
	"fmt"
	"sync"
)

//go:generate go run ../gen/main.go command-classes -c gen.config.yaml -o .
//go:generate go run ../gen/main.go parser -c gen.config.yaml -o ./command-classes.gen.go
//go:generate go run ../gen/main.go devices -c gen.config.yaml -o ./devices.gen.go

type (
	CommandClassID byte
	CommandID      byte

	Command interface {
		encoding.BinaryMarshaler
		encoding.BinaryUnmarshaler

		CommandClassID() CommandClassID
		CommandID() CommandID
		CommandIDString() string
	}

	CommandFactory func() Command

	CommandIdentifier struct {
		CommandClass CommandClassID
		Command      CommandID
		Version      uint8
	}
)

var (
	ErrNotRegistered    = errors.New("No factory exists for the specified command class")
	ErrPayloadUnderflow = errors.New("Payload length underflow")
)

var (
	factoriesMu sync.Mutex
	factories   = make(map[CommandIdentifier]CommandFactory)
)

func Register(identifier CommandIdentifier, factory CommandFactory) {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()

	if factory == nil {
		panic("commandclass: Register factory is nil")
	}

	if _, ok := factories[identifier]; ok {
		panic(fmt.Sprintf(
			"commandclass: Register called twice (cc: %x; c: %x; v: %d)\n",
			identifier.CommandClass,
			identifier.Command,
			identifier.Version,
		))
	}

	factories[identifier] = factory
}

func Parse(version uint8, payload []byte) (Command, error) {
	if payload == nil || len(payload) < 2 {
		return nil, ErrPayloadUnderflow
	}

	identifier := CommandIdentifier{
		CommandClass: CommandClassID(payload[0]),
		Command:      CommandID(payload[1]),
		Version:      version,
	}

	factoriesMu.Lock()
	factory, ok := factories[identifier]
	factoriesMu.Unlock()

	if !ok {
		return nil, ErrNotRegistered
	}

	command := factory()

	if err := command.UnmarshalBinary(payload); err != nil {
		return nil, err
	}

	return command, nil
}

// for tests
func unregisterAllFactories() {
	factoriesMu.Lock()
	defer factoriesMu.Unlock()
	factories = make(map[CommandIdentifier]CommandFactory)
}
