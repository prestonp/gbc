.PHONY: build
build:
	go build -o bin/gbc ./cmd/gbc

.PHONY: test
test:
	go test ./...
