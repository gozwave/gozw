package cc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommand struct{}

func (testCommand) CommandClassID() CommandClassID    { return 0x01 }
func (testCommand) CommandID() CommandID              { return 0x02 }
func (testCommand) CommandIDString() string           { return "" }
func (testCommand) MarshalBinary() ([]byte, error)    { return []byte{}, nil }
func (testCommand) UnmarshalBinary(data []byte) error { return nil }
func testCommandFactory() Command                     { return &testCommand{} }

var errTest = errors.New("")

type errCommand struct{}

func (errCommand) CommandClassID() CommandClassID    { return 0x02 }
func (errCommand) CommandID() CommandID              { return 0x02 }
func (errCommand) CommandIDString() string           { return "" }
func (errCommand) MarshalBinary() ([]byte, error)    { return nil, errTest }
func (errCommand) UnmarshalBinary(data []byte) error { return errTest }
func errCommandFactory() Command                     { return &errCommand{} }

func TestRegisterHandler(t *testing.T) {
	unregisterAllFactories()

	id := CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}

	handler := CommandFactory(func() Command {
		return &testCommand{}
	})

	assert.Len(t, factories, 0)

	Register(id, handler)

	assert.Len(t, factories, 1)

	_, ok := factories[id]().(*testCommand)
	assert.True(t, ok)
}

func TestRegisterNilFactoryPanics(t *testing.T) {
	unregisterAllFactories()

	assert.Panics(t, func() {
		Register(CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}, nil)
	})
}

func TestRegisterDuplicatePanics(t *testing.T) {
	unregisterAllFactories()

	assert.Panics(t, func() {
		Register(CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}, testCommandFactory)
		Register(CommandIdentifier{CommandClassID(0x02), CommandID(0x02), 1}, errCommandFactory)
		Register(CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}, testCommandFactory)
	})
}

func TestParse(t *testing.T) {
	unregisterAllFactories()
	Register(CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}, testCommandFactory)
	Register(CommandIdentifier{CommandClassID(0x02), CommandID(0x02), 1}, errCommandFactory)

	cmd, err := Parse(1, []byte{0x01, 0x02})
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.EqualValues(t, 0x01, cmd.CommandClassID())
	assert.EqualValues(t, 0x02, cmd.CommandID())
}

func TestUnregisteredParseReturnsError(t *testing.T) {
	unregisterAllFactories()

	_, err := Parse(1, []byte{0x01, 0x02})
	assert.Error(t, err)
	assert.Equal(t, ErrNotRegistered, err)
}

func TestParseBadPayloadSize(t *testing.T) {
	unregisterAllFactories()
	Register(CommandIdentifier{CommandClassID(0x01), CommandID(0x02), 1}, testCommandFactory)
	Register(CommandIdentifier{CommandClassID(0x02), CommandID(0x02), 1}, errCommandFactory)

	var err error

	_, err = Parse(1, nil)
	assert.Error(t, err)
	assert.Equal(t, err, ErrPayloadUnderflow)

	_, err = Parse(1, []byte{})
	assert.Error(t, err)
	assert.Equal(t, err, ErrPayloadUnderflow)

	_, err = Parse(1, []byte{0x00})
	assert.Error(t, err)
	assert.Equal(t, err, ErrPayloadUnderflow)
}

func TestParseReturnsUnmarshalError(t *testing.T) {
	unregisterAllFactories()
	Register(CommandIdentifier{CommandClassID(0x02), CommandID(0x02), 1}, errCommandFactory)

	var err error

	_, err = Parse(1, []byte{0x02, 0x02, 0x01})
	assert.Error(t, err)
	assert.Equal(t, err, errTest)

}
