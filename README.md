# Noversia Platform

Plateforme d'intelligence financière personnelle.

## Nouveautés v3

- Parsing CSV réel côté API Go.
- Rapport d'import ligne par ligne.
- Validation des colonnes obligatoires.
- Validation des montants et dates.
- Détection basique des lignes invalides.
- Suppression du fichier `GIT_COMMIT_MESSAGE.md` de l'archive.

## Lancement local

```bash
cp .env.example .env
docker compose up -d
make api
make ai
```

## Test import CSV

```bash
curl -X POST http://localhost:8080/api/v1/transactions/import \
  -H "Content-Type: multipart/form-data" \
  -F "file=@samples/bank-transactions-sample.csv"
```
