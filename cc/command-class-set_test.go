package cc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandClassSet(t *testing.T) {
	s := CommandClassSet{}

	assert.EqualValues(t, []CommandClassID{}, s.ListAll())
	assert.True(t, s.AllVersionsReceived())

	s.Add(Security)
	assert.EqualValues(t, []CommandClassID{Security}, s.ListAll())
	assert.False(t, s.AllVersionsReceived())

	assert.True(t, s.Supports(Security))
	assert.False(t, s.Supports(Alarm))
	assert.False(t, s.IsSecure(Security))
	assert.False(t, s.IsSecure(Alarm))
	assert.EqualValues(t, 0, s.GetVersion(Security))
	assert.EqualValues(t, 0, s.GetVersion(Alarm))

	s.SetVersion(Security, 1)
	assert.EqualValues(t, 1, s.GetVersion(Security))

	s.SetSecure(Security, true)
	assert.True(t, s.IsSecure(Security))

	s.SetSecure(Alarm, true)
	assert.True(t, s.Supports(Alarm))
	assert.True(t, s.IsSecure(Alarm))
	assert.EqualValues(t, 0, s.GetVersion(Alarm))
	assert.EqualValues(t, []CommandClassID{Security, Alarm}, s.ListAll())

	s.SetVersion(Association, 3)
	assert.True(t, s.Supports(Association))
	assert.False(t, s.IsSecure(Association))
	assert.EqualValues(t, 3, s.GetVersion(Association))

	assert.Contains(t, s.ListBySecureStatus(true), Security)
	assert.Contains(t, s.ListBySecureStatus(true), Alarm)
	assert.NotContains(t, s.ListBySecureStatus(true), Association)

	assert.NotContains(t, s.ListBySecureStatus(false), Security)
	assert.NotContains(t, s.ListBySecureStatus(false), Alarm)
	assert.Contains(t, s.ListBySecureStatus(false), Association)

	s.SetVersion(Alarm, 1)
	assert.True(t, s.AllVersionsReceived())

}
