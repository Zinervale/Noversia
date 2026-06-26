# US-004 — Persister les transactions importées

## Statut

Réalisé en v4.

## Critères d'acceptation

- Les transactions valides sont écrites en base.
- Les lignes invalides sont historisées.
- Un batch d'import est créé.
- Un rapport d'import est consultable.
- Les doublons simples ne sont pas réinsérés.
