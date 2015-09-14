package util

import (
	"testing"

	"github.com/helioslabs/gozw/cc"
	"github.com/stretchr/testify/assert"
)

func TestCommandClassSet(t *testing.T) {
	s := CommandClassSet{}

	assert.EqualValues(t, []cc.CommandClassID{}, s.ListAll())
	assert.True(t, s.AllVersionsReceived())

	s.Add(cc.Security)
	assert.EqualValues(t, []cc.CommandClassID{cc.Security}, s.ListAll())
	assert.False(t, s.AllVersionsReceived())

	assert.True(t, s.Supports(cc.Security))
	assert.False(t, s.Supports(cc.Alarm))
	assert.False(t, s.IsSecure(cc.Security))
	assert.False(t, s.IsSecure(cc.Alarm))
	assert.EqualValues(t, 0, s.GetVersion(cc.Security))
	assert.EqualValues(t, 0, s.GetVersion(cc.Alarm))

	s.SetVersion(cc.Security, 1)
	assert.EqualValues(t, 1, s.GetVersion(cc.Security))

	s.SetSecure(cc.Security, true)
	assert.True(t, s.IsSecure(cc.Security))

	s.SetSecure(cc.Alarm, true)
	assert.True(t, s.Supports(cc.Alarm))
	assert.True(t, s.IsSecure(cc.Alarm))
	assert.EqualValues(t, 0, s.GetVersion(cc.Alarm))
	assert.EqualValues(t, []cc.CommandClassID{cc.Security, cc.Alarm}, s.ListAll())

	s.SetVersion(cc.Association, 3)
	assert.True(t, s.Supports(cc.Association))
	assert.False(t, s.IsSecure(cc.Association))
	assert.EqualValues(t, 3, s.GetVersion(cc.Association))

	assert.EqualValues(t, []cc.CommandClassID{cc.Security, cc.Alarm}, s.ListBySecureStatus(true))
	assert.EqualValues(t, []cc.CommandClassID{cc.Association}, s.ListBySecureStatus(false))

	s.SetVersion(cc.Alarm, 1)
	assert.True(t, s.AllVersionsReceived())

}
