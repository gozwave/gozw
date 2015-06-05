package portal

type ZIPPacket struct {
	Payload []byte
}

func (zip *ZIPPacket) Marshal() []byte {
	return zip.Payload
}
