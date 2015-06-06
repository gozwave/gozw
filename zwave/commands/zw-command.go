package commands

type ZwFunction interface {
	Marshal() []byte
	Unmarshal() []byte
}
