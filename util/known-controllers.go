package util

// KnownUsbDevices returns a list of know USB devices
func KnownUsbDevices() map[string]string {
	return map[string]string{
		"0658:0200": "Aeotec Z-Stick Gen5 (ZW090-C)",
	}
}
