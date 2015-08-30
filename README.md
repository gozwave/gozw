# gozw

Golang Z/IP Portal

## Getting Started

1. You need a Z-Wave USB bridge controller with Serial API >= 4.05
1. With the controller plugged in, if you do not see a device in `/dev` named
   `usbserial*` or `usbmodem*` (or similar), you need to install the FTDI VCP
   driver for your platform
1. OPTIONAL: Install [`interceptty`](http://www.suspectclass.com/sgifford/interceptty/),
   which is an extremely useful serial port / tty proxy (it will allow you to see
   the raw data being transmitted to/from the controller).
  - Installation is as simple as `./configure && make && sudo cp interceptty /usr/local/bin`
  - To run: `interceptty -s 'ispeed 115200 ospeed 115200' /dev/<serialdevice> /tmp/<serialdevice>`
  - Be sure to update your config.yaml to point to `/tmp/<serialdevice>`
1. `go get -u github.com/helioslabs/zwgen`
1. `make install-deps` in the zwgen folder root
1. `make install` in the zwgen folder root
1. `make build` in the zwgen folder
1. `go get -u github.com/helioslabs/proto`
1. `go get ./...` in the gozw folder root
1. `go generate ./...` in the gozw folder root
1. `go run cmd/portald/main.go`
1. `go run cmd/gatewayd/main.go`

### Resources

- http://pepper1.net/zwavedb/ Random Z-Wave device library (lots of technical info)
