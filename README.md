# Noversia Platform

## v7 — Manual Category Correction

Cette version ajoute la correction manuelle de catégorie et l'historique d'enrichissement.

## Nouveautés

- Endpoint `PATCH /api/v1/transactions/{id}/category`.
- Table `transaction_enrichments`.
- Historique des corrections de catégorie.
- Détection de suggestion de règle après correction.
- Endpoint `GET /api/v1/rule-suggestions`.
- Documentation produit/technique mise à jour.

## Lancement

```bash
cp .env.example .env
docker compose up -d
make api
```

## Exemple correction catégorie

```bash
curl -X PATCH http://localhost:8080/api/v1/transactions/<id>/category \
  -H "Content-Type: application/json" \
  -d '{"categoryId":"<category-id>","reason":"manual_correction"}'
```
