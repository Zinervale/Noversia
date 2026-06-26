# Changelog

## v4 — Persist CSV Imports

### Ajouté

- Connexion PostgreSQL depuis le backend Go.
- Tables `import_batches` et `import_rows`.
- Persistance des transactions valides.
- Hash de transaction pour limiter les doublons.
- Endpoint de consultation d'un import.
- Lecture des transactions depuis PostgreSQL.

### Modifié

- `GET /transactions` ne retourne plus uniquement des données mockées.
- OpenAPI passe en v0.5.

### Prochaine étape

- Refactoriser le handler en couche service/repository.
- Ajouter des tests unitaires sur le parser CSV.
- Ajouter la catégorisation par règles.
