.PHONY: build test

build:
	echo "Building..."
	go build -o bin/coding-standards-rollout-tool main.go

test:
	go fmt ./...; \
	go build ./...; \
	go test -v ./...; \
	echo "\n\n";
