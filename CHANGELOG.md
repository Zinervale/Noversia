# Changelog

## v5 — Transactions Module Refactor

### Ajouté

- Parser CSV isolé dans le module transactions.
- Repository PostgreSQL dédié.
- Service d'import dédié.
- Handler HTTP dédié.
- Tests unitaires du parser CSV.
- Documentation de dette technique réduite.

### Modifié

- `cmd/api/main.go` ne contient plus la logique transaction.
- OpenAPI passe en v0.6.

### Prochaine étape

- Catégorisation automatique par règles.
- Tables `categorization_rules` et `transaction_enrichments`.
