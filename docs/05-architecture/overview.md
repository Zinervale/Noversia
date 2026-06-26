# Architecture globale

```text
Client -> API Go -> PostgreSQL
               -> AI Service Python
               -> Redis / Neo4j futurs
```

En v4, PostgreSQL devient la source de vérité pour les transactions importées.
