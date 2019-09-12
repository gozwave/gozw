.PHONY: bin

generate:
	@mkdir -p bin
	cd gen && go-bindata data/ templates/
	go build -o bin/gen ./gen
	./bin/gen devices -output=cc/devices_gen.go -config=cc/gen.config.yaml
	./bin/gen parser -output=cc/command_classes_gen.go -config=cc/gen.config.yaml
	./bin/gen command-classes -output=cc -config=cc/gen.config.yaml

bin:
	go build -o bin/basic ./examples/basic
