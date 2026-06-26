# Changelog

## v6 — Rule-Based Categorization

### Ajouté

- Moteur de catégorisation déterministe.
- Table `categorization_rules`.
- Catégories par défaut.
- Application des règles pendant l'import CSV.
- Endpoints de lecture/création des règles.
- Tests unitaires sur le matching de règles.

### Modifié

- Les transactions importées peuvent désormais recevoir une catégorie automatiquement.
- OpenAPI passe en v0.7.

### Prochaine étape

- Correction manuelle d'une catégorie.
- Création automatique d'une règle après corrections répétées.
- Enrichissement marchand.
