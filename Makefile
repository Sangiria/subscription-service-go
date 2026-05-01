PHONY: run build up down

run:
	go run cmd/server/main.go
build:
	docker compose up --build
up:
	docker compose up
down:
	docker compose down