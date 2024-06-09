# go-dating-app

## Requirements

- Docker

## Structure

- **.devcontainer**: dev container assets.
- **api**: handle request.
    - **rest**: request from rest api.
- **app**: app core function.
    - **dto**: data transfer object, transform request/response.
    - **entity**: represent domain object.
    - **repository**: manage data storage.
    - **service**: glue, a place to coordinate business/application flows.
- **cmd**: main file, to compose a binary.
- **common**: common functions collection.
    - **password**: password hashing
    - **validation**: var/struct validation.
- **config**: map environment variables into a struct.
- **database**: database related. i.e. migrations.
- **docs**: documentation. i.e. api specification.
- **scripts**: script for build/deploy.
- **storage**: data storage.

## Run

### Option 1: Dev container

**Requirements:**

- Dev container (VS code/Jetbrains)

**Steps:**

1. Copy `.env.example` on root directory to `.env`.
2. Adjust `.env` values when you make change to docker-compose.yml within `.devcontainer` directory.
3. Run dev container. It will take about a minute to set up.
4. Run migration with command `make migrate`
5. Run app with command `make run`

**Please note:**

- Services exposed via docker port (hard-coded):
    - App available on `:13000`
    - DB available on `:13306`
- If you replace `APP_ADDRESS` with something other than `:3000` then you need to adjust `.devcontainer` settings.

### Option 2: Local machine

**Requirements:**

- Go 1.22
- MySQL 8+ (app tested with v8.4)
- Make

**Steps:**

1. Prepare mysql server and a database to use.
2. Copy `.env.example` on root directory to `.env`.
3. Adjust the `.env` values, especially the database section.
4. Run migration with command `make migrate`
5. Run app with command `make run`

**Please note:**

- App port will follow `.env` setting.

## Deploy

## Option 1: Build binary

**Requirements:**

- Go 1.22
- Make

**Steps:**

1. Run `make build-binary`
2. It wil produce two binaries;
    - `dating`: main app. Run with `./dating`.
    - `migration`: database migration tools. Run with `./migration`.
    - Both binaries can detect environment variables or `.env` file.

## Option 2: Build image

**Requirements:**

- Docker
- Make

**Steps:**

1. Run `make build-image`. It will run multi-stage build.
2. It will produce single image named `datingapp`.
3. Inside there are two binaries, `app` and `migration`.
4. Image based on alpine image. Can do `docker exec -it CONTAINER_ID /bin/sh` to get into the container shell.
5. `migration` binary available on root directory. Run with `./migration`.
6. Make sure to pass environment variables as contained in the `.env` file.

## Tests

Test will produce a test coverage report, coverage.html, on root directory.

### Unit test

```shell
make test
```

### Integration test

```shell
# Only available with dev container.
# Will run unit + integration
make test-devc
```

Preview: [coverage report](https://raw.githack.com/alfi-repo/go-dating-app2/main/coverage.html)

## Dev container

- Based on official go 1.22 docker image.
- Alpine version.
- Installed tools:
    - golangci-lint v1.59 (config: ./.golangci.yml)
    - goose v3.20
    - mockery v2.43.2 (config: ./.mockery.yaml)
- Compatible with vscode/jetbrains remote development.
- Several vscode plugins have been added for ease of development.

## API Documentation (OpenAPI 3.1)

Please refer to `/docs` directory.  
Preview: [api doc](https://elements-demo.stoplight.io/?spec=https://raw.githubusercontent.com/alfi-repo/go-dating-app2/main/docs/openapi.yaml)
