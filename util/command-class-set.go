package util

import (
	"encoding/gob"

	"github.com/helioslabs/gozw/cc"
)

func init() {
	gob.Register(CommandClassSupport{})
	gob.Register(CommandClassSet{})
}

type CommandClassSupport struct {
	CommandClass cc.CommandClassID
	Secure       bool
	Version      uint8
}

type CommandClassSet map[cc.CommandClassID]*CommandClassSupport

func (s CommandClassSet) Supports(id cc.CommandClassID) bool {
	_, ok := s[id]
	return ok
}

func (s CommandClassSet) IsSecure(id cc.CommandClassID) bool {
	if c, ok := s[id]; ok {
		return c.Secure
	} else {
		return false
	}
}

func (s CommandClassSet) ListAll() []cc.CommandClassID {
	list := make([]cc.CommandClassID, 0)
	for id := range s {
		list = append(list, id)
	}
	return list
}

func (s CommandClassSet) ListBySecureStatus(secure bool) []cc.CommandClassID {
	list := make([]cc.CommandClassID, 0)
	for id, c := range s {
		if c.Secure == secure {
			list = append(list, id)
		}
	}
	return list
}

func (s CommandClassSet) GetVersion(id cc.CommandClassID) uint8 {
	if c, ok := s[id]; ok {
		return c.Version
	} else {
		return 0
	}
}

func (s CommandClassSet) Add(id cc.CommandClassID) {
	_, ok := s[id]
	if !ok {
		s[id] = &CommandClassSupport{
			CommandClass: id,
		}
	}
}

func (s CommandClassSet) SetSecure(id cc.CommandClassID, secure bool) {
	if c, ok := s[id]; ok {
		c.Secure = secure
	} else {
		s[id] = &CommandClassSupport{
			CommandClass: id,
			Secure:       secure,
		}
	}
}

func (s CommandClassSet) SetVersion(id cc.CommandClassID, version uint8) {
	if c, ok := s[id]; ok {
		c.Version = version
	} else {
		s[id] = &CommandClassSupport{
			CommandClass: id,
			Version:      version,
		}
	}
}

func (s CommandClassSet) AllVersionsReceived() bool {
	for _, c := range s {
		if c.Version == 0 {
			return false
		}
	}

	return true
}
