ENTRYPOINT=cmd/shortener/main.go

tidy:
	@go mod tidy

run: tidy
	@go run ${ENTRYPOINT}

sch:
	staticcheck ./...

test-all:
	go test -v ./...

