.PHONY: test

test:
	env GOMAXPROCS=8 go test -cover ./zwave/...
