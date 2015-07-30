package application

import "github.com/bjyoungblood/gozw/zwave/serial-api"

type ApplicationLayer struct {
	serialApi serialapi.ISerialAPILayer
}

	appLayer := &ApplicationLayer{
func NewApplicationLayer(serialApi serialapi.ISerialAPILayer) *ApplicationLayer {
		serialApi: serialApi,
	}

	return appLayer
}

func (*ApplicationLayer) SendData() {

}

func (*ApplicationLayer) applicationCommandHandler() {

}
