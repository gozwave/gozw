package util

import (
	"encoding/gob"

	"github.com/helioslabs/gozw/zwave/command-class"
)

func init() {
	gob.Register(CommandClassSupport{})
	gob.Register(CommandClassSet{})
}

type CommandClassSupport struct {
	CommandClass commandclass.ID
	Secure       bool
	Version      uint8
}

type CommandClassSet map[commandclass.ID]*CommandClassSupport

func (s CommandClassSet) Supports(id commandclass.ID) bool {
	_, ok := s[id]
	return ok
}

func (s CommandClassSet) IsSecure(id commandclass.ID) bool {
	if c, ok := s[id]; ok {
		return c.Secure
	} else {
		return false
	}
}

func (s CommandClassSet) ListAll() []commandclass.ID {
	list := make([]commandclass.ID, 0)
	for id := range s {
		list = append(list, id)
	}
	return list
}

func (s CommandClassSet) ListBySecureStatus(secure bool) []commandclass.ID {
	list := make([]commandclass.ID, 0)
	for id, c := range s {
		if c.Secure == secure {
			list = append(list, id)
		}
	}
	return list
}

func (s CommandClassSet) GetVersion(id commandclass.ID) uint8 {
	if c, ok := s[id]; ok {
		return c.Version
	} else {
		return 0
	}
}

func (s CommandClassSet) Add(id commandclass.ID) {
	_, ok := s[id]
	if !ok {
		s[id] = &CommandClassSupport{
			CommandClass: id,
		}
	}
}

func (s CommandClassSet) SetSecure(id commandclass.ID, secure bool) {
	if c, ok := s[id]; ok {
		c.Secure = secure
	} else {
		s[id] = &CommandClassSupport{
			CommandClass: id,
			Secure:       secure,
		}
	}
}

func (s CommandClassSet) SetVersion(id commandclass.ID, version uint8) {
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
