# ADR 0005 — Persister les imports CSV dans PostgreSQL

## Statut

Accepté

## Décision

Les imports CSV sont stockés dans PostgreSQL via `import_batches`, `import_rows` et `transactions`.

## Justification

- Traçabilité.
- Audit.
- Débogage facile.
- Préparation RGPD/export/suppression.
