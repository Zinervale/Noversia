# Noversia Platform

## v6 — Rule-Based Categorization

Cette version ajoute la première brique d'intelligence déterministe : la catégorisation automatique par règles.

## Nouveautés

- Tables `categorization_rules` et catégories seedées.
- Catégorisation automatique à l'import CSV.
- Endpoint `GET /api/v1/categories`.
- Endpoint `GET /api/v1/categorization-rules`.
- Endpoint `POST /api/v1/categorization-rules`.
- Tests unitaires du moteur de catégorisation.
- Documentation fonctionnelle et technique.

## Lancement

```bash
cp .env.example .env
docker compose up -d
make api
make test
```
