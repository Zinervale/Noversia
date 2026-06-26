# Module Transactions

## Responsabilités

- Lire les transactions.
- Importer un CSV.
- Valider les lignes.
- Persister les transactions.
- Historiser les lignes invalides.
- Éviter les doublons simples.

## Interfaces

### Parser

Transforme un fichier CSV en `ImportReport`.

### Service

Orchestre le parsing et la persistance.

### Repository

Accède à PostgreSQL.

### Handler

Expose les routes HTTP.
