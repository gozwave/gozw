.PHONY: test

generate:
	go generate ./cc

interceptty:
	interceptty -s 'ispeed 115200 ospeed 115200' /dev/cu.usbmodem1411 /tmp/usbmodem

install: generate

lint:
	gometalinter -e ".gen.go" ./...

test:
	gometalinter -D golint -D errcheck -D gocyclo -D dupl ./...
	env GOMAXPROCS=8 go test -cover ./zwave/...

install-tools:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update
