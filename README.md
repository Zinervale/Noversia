# Noversia Platform

Plateforme d'intelligence financière personnelle.

## Nouveautés v2

- Ajout du `CHANGELOG.md`
- Ajout du `GIT_COMMIT_MESSAGE.md`
- Ajout d'un exemple CSV bancaire
- Ajout d'un endpoint d'import CSV simulé
- Ajout d'une documentation fonctionnelle de l'import
- Ajout d'une structure métier Transaction Import

## Stack

- Backend Core : Go
- Service IA : Python / FastAPI
- Base métier : PostgreSQL
- Cache : Redis
- Graphe : Neo4j
- Documentation : MkDocs
- API : REST + OpenAPI

## Lancement local

```bash
cp .env.example .env
docker compose up -d
make api
make ai
```

## Test rapide API

```bash
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/transactions
curl -X POST http://localhost:8080/api/v1/transactions/import \
  -H "Content-Type: multipart/form-data" \
  -F "file=@samples/bank-transactions-sample.csv"
```
