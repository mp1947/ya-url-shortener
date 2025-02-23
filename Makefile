ENTRYPOINT=cmd/shortener/main.go
GO_VERSION=1.22.12
APP_NAME=shortener

tidy:
	@go mod tidy -go=${GO_VERSION}

build: tidy
	@go build -o ./bin/${APP_NAME} ${ENTRYPOINT}

run-debug: build
	@GIN_MODE=debug ./bin/${APP_NAME} ${ARGS}

run-release: build
	@GIN_MODE=release ./bin/${APP_NAME} ${ARGS}

check:
	staticcheck ./...
	go vet ./...

test-all:
	go test -v ./...

lint:
	golangci-lint run --issues-exit-code 1 --print-issued-lines=true  ./...

docker-test:
	docker buildx build . \
		--build-arg APP_PATH=${ENTRYPOINT}