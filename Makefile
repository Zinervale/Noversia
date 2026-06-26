.PHONY: up down api ai docs test

up:
	docker compose up -d

down:
	docker compose down

api:
	cd services/api && go run ./cmd/api

ai:
	cd services/ai && python -m uvicorn app.main:app --reload --port 8000

docs:
	mkdocs serve

test:
	cd services/api && go test ./...
