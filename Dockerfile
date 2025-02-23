# syntax=docker/dockerfile:1.4
FROM golang:1.22.12 AS builder

ARG APP_PATH

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o ./bin/app ${APP_PATH}


FROM alpine:3.21.2

WORKDIR /

COPY --from=builder /build/bin/app .

ENTRYPOINT [ "/app" ]