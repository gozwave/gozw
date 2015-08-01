package protocol

const (
	TransmitOptionAck       byte = 0x01
	TransmitOptionLowPower       = 0x02
	TransmitOptionAutoRoute      = 0x04
	TransmitOptionNoRoute        = 0x10
	TransmitOptionExplore        = 0x20
)

const (
	TransmitCompleteOk      byte = 0x00
	TransmitCompleteNoAck        = 0x01
	TransmitCompleteFail         = 0x02
	TransmitRoutingNotIdle       = 0x03
	TransmitCompleteNoRoute      = 0x04
)
