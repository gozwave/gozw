build:
	go generate ./...
	go build ./...

install: build
	go install .

install-deps:
	go get -u github.com/jteeuwen/go-bindata/...
	glide install
