default: test

.PHONY: test
test:
	go test ./...
	go run ./cmd/bfc -r -o /tmp/x ./_examples/hello.bf | grep 'Hello World!'

