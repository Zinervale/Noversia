# Standards API

API REST versionnée : `/api/v1`.

## Erreur standard

```json
{
  "error": {
    "code": "TRANSACTION_NOT_FOUND",
    "message": "Transaction introuvable",
    "correlationId": "..."
  }
}
```

## Sécurité
OAuth2 / OIDC, JWT court terme, refresh token rotatif, scopes API, rate limiting, audit.
