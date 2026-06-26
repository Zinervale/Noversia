# Changelog

## v2 — Transaction Import Foundation

### Ajouté

- `GIT_COMMIT_MESSAGE.md` avec le message Git prêt à copier.
- Exemple `samples/bank-transactions-sample.csv`.
- Endpoint `POST /api/v1/transactions/import`.
- Structure métier initiale pour l'import CSV.
- Documentation de l'import bancaire.
- OpenAPI v0.3.

### Modifié

- README enrichi avec commandes de test.
- Roadmap mise à jour.

### Note

L'import CSV est encore simulé côté API. La prochaine étape sera le parsing réel, la validation ligne par ligne et l'écriture en base.
