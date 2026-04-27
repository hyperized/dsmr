.PHONY: build test test-coverage test-integration lint clean

build:
	go build ./...

test:
	go test ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

test-integration:
	go test -tags integration ./...

lint:
	golangci-lint run ./...

clean:
	rm -f coverage.out
