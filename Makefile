-include .env

ifeq ($(TARGET_OS),windows)
	EXT := .exe
	RM	:= del /Q
	BIN	:= bin\\
else
	EXT :=
	RM	:= rm -f
	BIN := bin/
endif

.PHONY: run migrate-up migrate-down migrate-version migrate-new

run:
	go run ./cmd/api

build:
	go build -o $(BIN)api$(EXT) ./cmd/api
	go build -o $(BIN)migrate$(EXT) ./cmd/migrate

clean:
	$(RM) $(BIN)api$(EXT) $(BIN)migrate$(EXT)

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down 1

migrate-version:
	go run ./cmd/migrate version

# usage: make migrate-new name=add_orders
migrate-new:
	go run ./cmd/migrate create $(name)