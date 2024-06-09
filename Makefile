# build docker image.
.PHONY: build-image
build-image:
	docker build -t datingapp -f scripts/build/Dockerfile .

# build binary.
.PHONY: build-binary
build-binary:
	go build -o ./dating ./cmd/dating/main.go
	go build -o ./migration ./cmd/dating/main.go

# generate mocks.
.PHONY: gen-mock
gen-mock:
	mockery --all

# run database migrations, read db connection from os env or .env.
.PHONY: migrate
migrate:
	go run cmd/migration/main.go

# run the app.
.PHONY: run
run:
	go run cmd/dating/main.go

# run unit test.
.PHONY: test
test:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o coverage.html

# run unit & integration with dev container.
.PHONY: test-devc
test-devc:
	TEST_INTEGRATION=devcontainer go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o coverage.html

# run unit test with report & verbose.
.PHONY: test-verbose
test-verbose:
	go test -v -coverprofile cover.out ./...
	go tool cover -html=cover.out -o coverage.html
