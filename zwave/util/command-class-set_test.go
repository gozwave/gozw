package util

import (
	"testing"

	"github.com/helioslabs/gozw/zwave/command-class"
	"github.com/stretchr/testify/assert"
)

func TestCommandClassSet(t *testing.T) {
	s := CommandClassSet{}

	assert.EqualValues(t, []commandclass.ID{}, s.ListAll())
	assert.True(t, s.AllVersionsReceived())

	s.Add(commandclass.Security)
	assert.EqualValues(t, []commandclass.ID{commandclass.Security}, s.ListAll())
	assert.False(t, s.AllVersionsReceived())

	assert.True(t, s.Supports(commandclass.Security))
	assert.False(t, s.Supports(commandclass.Alarm))
	assert.False(t, s.IsSecure(commandclass.Security))
	assert.False(t, s.IsSecure(commandclass.Alarm))
	assert.EqualValues(t, 0, s.GetVersion(commandclass.Security))
	assert.EqualValues(t, 0, s.GetVersion(commandclass.Alarm))

	s.SetVersion(commandclass.Security, 1)
	assert.EqualValues(t, 1, s.GetVersion(commandclass.Security))

	s.SetSecure(commandclass.Security, true)
	assert.True(t, s.IsSecure(commandclass.Security))

	s.SetSecure(commandclass.Alarm, true)
	assert.True(t, s.Supports(commandclass.Alarm))
	assert.True(t, s.IsSecure(commandclass.Alarm))
	assert.EqualValues(t, 0, s.GetVersion(commandclass.Alarm))
	assert.EqualValues(t, []commandclass.ID{commandclass.Security, commandclass.Alarm}, s.ListAll())

	s.SetVersion(commandclass.Association, 3)
	assert.True(t, s.Supports(commandclass.Association))
	assert.False(t, s.IsSecure(commandclass.Association))
	assert.EqualValues(t, 3, s.GetVersion(commandclass.Association))

	assert.EqualValues(t, []commandclass.ID{commandclass.Security, commandclass.Alarm}, s.ListBySecureStatus(true))
	assert.EqualValues(t, []commandclass.ID{commandclass.Association}, s.ListBySecureStatus(false))

	s.SetVersion(commandclass.Alarm, 1)
	assert.True(t, s.AllVersionsReceived())

}
