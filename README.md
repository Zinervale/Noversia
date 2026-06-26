# Noversia Platform

Plateforme d'intelligence financière personnelle.

## Objectif

Construire un moteur de décision financière capable de :
- centraliser des données financières ;
- analyser les transactions ;
- produire des recommandations explicables ;
- simuler des décisions de vie ;
- dialoguer avec l'utilisateur via une IA maîtrisée.

## Stack v1

- Backend Core : Go
- Service IA : Python / FastAPI
- Base métier : PostgreSQL
- Cache : Redis
- Graphe : Neo4j
- Documentation : MkDocs
- API : REST + OpenAPI

## Lancement local

```bash
cp .env.example .env
docker compose up -d
make api
make ai
```

## Structure

```text
apps/          Applications clientes
services/      Services backend et IA
packages/      Librairies partagées
docs/          Documentation projet
adr/           Décisions d'architecture
openapi/       Contrats API
backlog/       Backlog produit
infra/         Infrastructure
scripts/       Scripts de développement
```
