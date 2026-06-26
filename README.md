# Noversia Platform

## v8 — Merchant Detection and Rule Suggestions

Cette version ajoute les marchands et transforme les suggestions en règles.

## Nouveautés

- Détection simple du marchand à partir du libellé.
- Normalisation du nom marchand.
- Création automatique/réutilisation du marchand pendant l'import.
- Ajout du `merchant_id` sur les transactions importées.
- Endpoint `POST /api/v1/rule-suggestions/apply`.
- Documentation marchands et règles enrichie.

## Lancement

```bash
cp .env.example .env
docker compose up -d
make api
```
