# Backend Go

## Structure

```text
services/api/
  cmd/api
  internal/config
  internal/http
  internal/modules
  internal/platform
```

## Principes

- HTTP standard library pour démarrer simplement
- Modules métier isolés
- Pas de dépendance lourde prématurée
- Migration SQL simple
- OpenAPI comme contrat externe
