.PHONY: docker-build docker-up docker-down run-local test

run-local:
	go run cmd/server/main.go
docker-build:
	docker compose up --build
docker-up:
	docker compose up
docker-down:
	docker compose down
test:
	go test -v ./internal/handlers