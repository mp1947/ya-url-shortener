ENTRYPOINT=cmd/shortener/main.go
GO_VERSION=1.22.12
APP_NAME=shortener

tidy:
	@go mod tidy -go=${GO_VERSION}

build: tidy
	@go build -o ./bin/${APP_NAME} ${ENTRYPOINT}

run: build
	@GIN_MODE=release ./bin/${APP_NAME} ${ARGS}

run-debug: build
	@GIN_MODE=debug ./bin/${APP_NAME} ${ARGS}

check-code:
	staticcheck ./...
	go vet ./...
	golangci-lint run --issues-exit-code 1  ./...

test-all:
	go test -v ./...

up:
	docker compose up -d

down:
	docker compose down

update-workflows:
	git fetch template && git checkout template/main .github
