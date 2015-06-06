package functions

type ZwFunction interface {
	Marshal() []byte
	Unmarshal() []byte
}
