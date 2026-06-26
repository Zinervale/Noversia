# ADR 0006 — Refactoriser le module Transactions

## Statut

Accepté

## Contexte

La logique d'import était dans `main.go`, ce qui rendait le code difficile à tester et à maintenir.

## Décision

Créer un module transactions en quatre couches :

- handler
- service
- repository
- parser

## Conséquences

- Plus de fichiers.
- Architecture plus propre.
- Tests plus simples.
- Préparation à l'import multi-banques.
