# Noversia Platform

Plateforme d'intelligence financière personnelle.

## Nouveautés v4

- Connexion PostgreSQL réelle côté API Go.
- Création des tables `import_batches` et `import_rows`.
- Persistance des transactions valides importées par CSV.
- Détection simple des doublons via empreinte `source_hash`.
- Endpoint `GET /api/v1/imports/{id}` pour consulter un rapport d'import.
- Endpoint `GET /api/v1/transactions` connecté à PostgreSQL.
- OpenAPI v0.5.

## Lancement local

```bash
cp .env.example .env
docker compose up -d
make api
make ai
```

## Test

```bash
curl -X POST http://localhost:8080/api/v1/transactions/import \
  -H "Content-Type: multipart/form-data" \
  -F "file=@samples/bank-transactions-sample.csv"

curl http://localhost:8080/api/v1/transactions
```
