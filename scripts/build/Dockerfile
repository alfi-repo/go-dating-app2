# syntax=docker/dockerfile:1
ARG GO_VERSION=1.22

## Build
FROM golang:${GO_VERSION} AS build
WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/app ./cmd/dating/main.go
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /go/bin/migration ./cmd/migration/main.go

## Deploy
FROM alpine:latest
WORKDIR /
COPY --from=build /go/bin/app /
COPY --from=build /go/bin/migration /
EXPOSE 3000
ENTRYPOINT ["/app"]
