package cc

import "encoding/gob"

func init() {
	gob.Register(CommandClassSupport{})
	gob.Register(CommandClassSet{})
}

type CommandClassSupport struct {
	CommandClass CommandClassID
	Secure       bool
	Version      uint8
}

type CommandClassSet map[CommandClassID]*CommandClassSupport

func (s CommandClassSet) Supports(id CommandClassID) bool {
	_, ok := s[id]
	return ok
}

func (s CommandClassSet) IsSecure(id CommandClassID) bool {
	if c, ok := s[id]; ok {
		return c.Secure
	} else {
		return false
	}
}

func (s CommandClassSet) ListAll() []CommandClassID {
	list := make([]CommandClassID, 0)
	for id := range s {
		list = append(list, id)
	}
	return list
}

func (s CommandClassSet) ListBySecureStatus(secure bool) []CommandClassID {
	list := make([]CommandClassID, 0)
	for id, c := range s {
		if c.Secure == secure {
			list = append(list, id)
		}
	}
	return list
}

func (s CommandClassSet) GetVersion(id CommandClassID) uint8 {
	if c, ok := s[id]; ok {
		return c.Version
	}

	return 0

}

func (s CommandClassSet) Add(id CommandClassID) {
	_, ok := s[id]
	if !ok {
		s[id] = &CommandClassSupport{
			CommandClass: id,
		}
	}
}

func (s CommandClassSet) SetSecure(id CommandClassID, secure bool) {
	if c, ok := s[id]; ok {
		c.Secure = secure
	} else {
		s[id] = &CommandClassSupport{
			CommandClass: id,
			Secure:       secure,
		}
	}
}

func (s CommandClassSet) SetVersion(id CommandClassID, version uint8) {
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
