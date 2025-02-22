ENTRYPOINT=cmd/shortener/main.go

tidy:
	@go mod tidy

run: tidy
	@go run ${ENTRYPOINT}

check:
	staticcheck ./...
	go vet ./...

test-all:
	go test -v ./...