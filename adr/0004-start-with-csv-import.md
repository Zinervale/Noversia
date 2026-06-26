# ADR 0004 — Démarrer par l'import CSV

## Statut

Accepté

## Contexte

La synchronisation bancaire automatique nécessite un fournisseur Open Banking, des tests de compatibilité bancaire et des contraintes réglementaires/commerciales.

## Décision

Le MVP démarre avec l'import CSV.

## Conséquences positives

- Démarrage rapide.
- Coût réduit.
- Moins de dépendances.
- Données réelles disponibles pour tester l'IA.

## Conséquences négatives

- Expérience moins fluide qu'une synchronisation bancaire automatique.
- Formats CSV variables selon les banques.
- Mapping colonnes à prévoir.
