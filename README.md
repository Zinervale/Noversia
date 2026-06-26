# Noversia Platform

## v5 — Transactions Module Refactor

Cette version refactorise le module transaction pour sortir la logique métier de `main.go`.

## Nouveautés

- `parser.go` : parsing et validation CSV.
- `repository.go` : accès PostgreSQL.
- `service.go` : orchestration import.
- `handler.go` : endpoints HTTP.
- Tests unitaires sur le parser CSV.
- Documentation architecture mise à jour.
- Préparation de la catégorisation par règles.

## Lancement

```bash
cp .env.example .env
docker compose up -d
make api
make ai
```

## Tests

```bash
make test
```
