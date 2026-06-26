# Changelog

## v3 — Real CSV Parsing

### Ajouté

- Parsing CSV réel côté Go.
- Validation des colonnes `date`, `label`, `amount`, `currency`.
- Rapport d'import ligne par ligne.
- Comptage des lignes valides et invalides.
- Exemple CSV avec une ligne invalide pour tester les erreurs.
- Documentation d'import CSV enrichie.

### Modifié

- L'archive ne contient plus `GIT_COMMIT_MESSAGE.md`.
- OpenAPI passe en v0.4.

### Prochaine étape

- Persistance PostgreSQL des transactions importées.
- Création d'une table `import_batches`.
- Détection de doublons.
