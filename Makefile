springpath-docker-plugin:
	go build $@
	go install $@

clean:
	go clean

test:
	go test ./...
