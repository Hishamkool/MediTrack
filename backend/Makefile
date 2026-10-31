# Run these from the backend/ folder, e.g. `make run`.

run:            ## start the API (also runs migrations + seeds admin)
	go run ./cmd/api

migrate-up:     ## apply all pending migrations
	go run ./cmd/migrate up

migrate-down:   ## roll back the last migration
	go run ./cmd/migrate down

migrate-version:
	go run ./cmd/migrate version

test:
	go test ./...

build:          ## produce a single binary at bin/api
	go build -o bin/api ./cmd/api

.PHONY: run migrate-up migrate-down migrate-version test build
