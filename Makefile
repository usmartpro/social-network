BIN := "./bin/social"
DOCKER_IMG="social-network:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

DB_CONN := "postgresql://postgres:postgres@localhost:5432/social?sslmode=disable"

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/social

run: build up
	#$(BIN)

version: build
	$(BIN) version

test:
	go test -race ./internal/... -count 100

lint:
	golangci-lint run ./...

migrate:
	goose --dir=migrations postgres ${DB_CONN} up

up:
	docker-compose up -d

down:
	docker-compose down
