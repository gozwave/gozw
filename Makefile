.PHONY: test

generate:
	go generate ./gen
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

cover:
	@echo Running coverage
	go get github.com/wadey/gocovmerge
	$(eval PKGS := $(shell go list ./... ))
	$(eval PKGS_DELIM := $(shell echo $(PKGS) | sed -e 's/ /,/g'))
	go list -f '{{if or (len .TestGoFiles) (len .XTestGoFiles)}}go test -test.v -test.timeout=120s -covermode=atomic -coverprofile={{.Name}}_{{len .Imports}}_{{len .Deps}}.coverprofile -coverpkg $(PKGS_DELIM) {{.ImportPath}}{{end}}' $(PKGS) | xargs -I {} bash -c {}
	gocovmerge `ls *.coverprofile` > coverage.txt
	rm *.coverprofile
