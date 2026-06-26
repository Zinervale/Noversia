.PHONY: docs up down

docs:
	mkdocs serve

up:
	docker compose up -d

down:
	docker compose down
