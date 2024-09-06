run:
	@go run cmd/main.go

CURRENT_DIR := $(shell pwd)
DB_URL := "postgres://postgres:azamat@localhost:5432/fitness?sslmode=disable"

migrate-up:
	migrate -path migrations/ -database $(DB_URL) up


migrate-down:
	migrate -path migrations/ -database  $(DB_URL) down

migrate-force:
	migrate -path migrations/ -database  $(DB_URL) force 1

sqlc-generate:
	@sqlc vet && sqlc generate