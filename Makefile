SHORTENER_ENTRYPOINT=cmd/shortener/main.go
MULTICHECKER_ENTRYPOINT=cmd/staticlint/main.go
GO_VERSION=1.24.4
SHORTENER_NAME=shortener
MULTICHECKER_NAME=multichecker
MOCKS_SOURCE=internal/repository/repository.go
MOCKS_DEST=internal/mocks/mock_repository.go

.PHONY: tidy build run run-debug check-code test bench mock

tidy:
	@go mod tidy -go=${GO_VERSION}

build: tidy mock
	@go build -o ./bin/${SHORTENER_NAME} ${SHORTENER_ENTRYPOINT}

build-multichecker: tidy
	@go build -o ./bin/${MULTICHECKER_NAME} ${MULTICHECKER_ENTRYPOINT}

run: build mock
	@GIN_MODE=release ./bin/${SHORTENER_NAME} ${ARGS}

run-debug: build up
	@GIN_MODE=debug ./bin/${SHORTENER_NAME} ${ARGS}

multichecker: build-multichecker
	go list ./... | grep -v mocks | xargs ./bin/${MULTICHECKER_NAME} -test=false

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
