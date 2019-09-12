package serialapi

// func (s *ZWaveSessionLayer) ApplicationNodeInformation(
// 	deviceOptions uint8,
// 	genericType uint8,
// 	specificType uint8,
// 	supportedCommandClasses []uint8,
// ) {
//
// 	s.execLock.Lock()
// 	defer s.execLock.Unlock()
//
// 	payload := []byte{
// 		FnApplicationNodeInformation,
// 		deviceOptions,
// 		genericType,
// 		specificType,
// 		uint8(len(supportedCommandClasses)),
// 	}
//
// 	payload = append(payload, supportedCommandClasses...)
//
// 	s.write(frame.NewRequestFrame(payload))
//
// }
