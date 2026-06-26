# Backend Go

## v4

Le backend Go se connecte réellement à PostgreSQL via `database/sql` et le driver `pgx`.

## Dette technique acceptée

La logique d'import reste dans `main.go` pour accélérer.

## Refactorisation v5

Extraire :

- parser CSV ;
- repository transaction ;
- service import ;
- handler HTTP.
