# Import CSV bancaire

## v5

Le flux d'import est maintenant découpé :

```mermaid
flowchart TD
    H[HTTP Handler] --> S[Import Service]
    S --> P[CSV Parser]
    S --> R[Transaction Repository]
    R --> DB[(PostgreSQL)]
```

## Bénéfices

- Testabilité.
- Lisibilité.
- Préparation à plusieurs formats bancaires.
- Préparation aux règles de catégorisation.
