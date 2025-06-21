ENTRYPOINT=cmd/shortener/main.go
GO_VERSION=1.24.4
APP_NAME=shortener
MOCKS_SOURCE=internal/repository/repository.go
MOCKS_DEST=internal/mocks/mock_repository.go

.PHONY: tidy build run run-debug check-code test bench mock

tidy:
	@go mod tidy -go=${GO_VERSION}

build: tidy mock
	@go build -o ./bin/${APP_NAME} ${ENTRYPOINT}

run: build up mock
	@GIN_MODE=release ./bin/${APP_NAME} ${ARGS}

run-debug: build up
	@GIN_MODE=debug ./bin/${APP_NAME} ${ARGS}

check-code:
	staticcheck ./...
	go vet ./...
	golangci-lint run  ./...

test: mock
	go test -v ./...

bench:
	go test -bench=. -benchmem -benchtime=10s -run=^Benchmark ./...

mock: tidy
	mockgen -source=${MOCKS_SOURCE} -destination=${MOCKS_DEST} -package=mocks

coverage:
	go test -covermode=count -coverprofile=coverage.out ./...
	grep -vE "mocks|repository/database|repository/inmemory" coverage.out > coverage.cleaned.out
	go tool cover -func=coverage.cleaned.out

up:
	docker compose up -d

down:
	docker compose down

update-workflows:
	git fetch template && git checkout template/main .github
