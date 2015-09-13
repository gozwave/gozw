package util

type ByteMarshaler []byte

func (b ByteMarshaler) MarshalBinary() ([]byte, error) {
	return b, nil
}
