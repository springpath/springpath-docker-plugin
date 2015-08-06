build:
	go build .
	go install .

clean:
	go clean

test:
	go test ./...

.PHONY: build clean test
