# Architecture globale

## Choix

Monolithe modulaire Go + service IA Python séparé.

```text
Client
  |
API Go
  |
PostgreSQL / Redis / Neo4j
  |
AI Service Python
  |
LLM Gateway
```

## Raison

Cette architecture permet de livrer vite sans s'enfermer dans des microservices prématurés.
