package zwave

const (
	FrameSOFData uint8 = 0x01
	FrameSOFAck  uint8 = 0x06
	FrameSOFNak  uint8 = 0x15
	FrameSOFCan  uint8 = 0x18
)

const (
	FrameTypeReq uint8 = 0x00
	FrameTypeRes uint8 = 0x01
)
