# US-005 — Refactoriser l'import transaction

## Statut

Réalisé en v5.

## Critères d'acceptation

- Le parsing est isolé.
- La persistance est isolée.
- Le handler HTTP ne porte pas la logique métier.
- Le parser possède des tests unitaires.
