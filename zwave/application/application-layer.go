package application

import "github.com/bjyoungblood/gozw/zwave/serial-api"

type ApplicationLayer struct {
	serialApi *serialapi.SerialAPILayer
}

func NewApplicationLayer(serialApi *serialapi.SerialAPILayer) *ApplicationLayer {
	appLayer := &ApplicationLayer{
		serialApi: serialApi,
	}

	return appLayer
}

func (*ApplicationLayer) SendData() {

}

func (*ApplicationLayer) applicationCommandHandler() {

}
