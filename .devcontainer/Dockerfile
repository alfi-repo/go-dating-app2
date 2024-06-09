ARG GO_VERSION=1.22

FROM golang:${GO_VERSION}-alpine

ARG GOLANGCILINT_VERSION=v1.59.0
ARG GOOSE_VERSION=v3.20.0
ARG MOCKERY_VERSION=2.43.2

RUN apk add --no-cache libxext libxrender libxtst libxi freetype procps gcompat musl libgcc libstdc++ bash curl make

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCILINT_VERSION} && \
    curl -sSfL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh -s ${GOOSE_VERSION} && \
    go install golang.org/x/tools/cmd/goimports@latest

COPY /scripts/mockery_install.sh /scripts/mockery_install.sh

RUN /scripts/mockery_install.sh ${MOCKERY_VERSION}
