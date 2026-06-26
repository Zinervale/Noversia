# Catégorisation automatique

## Objectif

Réduire le reclassement manuel mensuel.

## Approche v6

Le système utilise des règles déterministes :

- type de règle : `contains`
- champ analysé : libellé de transaction
- catégorie cible
- priorité
- score de confiance

## Exemples

| Motif | Catégorie |
|---|---|
| CARREFOUR | Courses |
| NETFLIX | Abonnements |
| SALAIRE | Revenus |
| TOTAL | Transport |

## Pourquoi commencer par des règles ?

- Prévisible.
- Explicable.
- Peu coûteux.
- Conforme au principe : l'IA n'invente pas.

L'IA viendra ensuite en complément lorsque les règles ne suffisent pas.
