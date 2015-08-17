.PHONY: test

generate:
	go generate ./zwave/command-class

clean-gen:
	-rm zwave/command-class/*.gen.go
	-rm zwave/command-class/**/*.gen.go
	-rmdir zwave/command-class/*

install: generate

lint:
	-go install ./...
	gometalinter -e ".gen.go" ./...

test:
	go install ./...
	gometalinter -D golint -D errcheck -D gocyclo -D dupl ./...
	env GOMAXPROCS=8 go test -cover ./zwave/...

install-tools:
	go get github.com/alecthomas/gometalinter
	gometalinter --install --update
