# Noversia Platform

Noversia Platform est une plateforme d'intelligence financière personnelle.

Objectif : construire un moteur de décision financière capable d'analyser les données bancaires, patrimoniales et comportementales d'un utilisateur pour produire des recommandations explicables, traçables et utiles.

## Produits

- **Noversia Financial** : application grand public.
- **Noversia Core** : moteur financier.
- **Noversia Intelligence** : moteur IA.
- **Noversia Decision Engine** : moteur de recommandations et simulations.
- **Noversia API** : API publique et privée.

## Structure

```text
apps/              Applications clientes
services/          Services backend
packages/          Bibliothèques partagées
docs/              Documentation produit et technique
adr/               Architecture Decision Records
openapi/           Contrats API
backlog/           Epics, features, user stories
diagrams/          Diagrammes C4, Mermaid, UML
infra/             Docker, Kubernetes, IaC
scripts/           Scripts de développement
```

## Décisions initiales

- Backend principal : Go
- Service IA : Python
- Base métier : PostgreSQL
- Cache : Redis
- Graphe financier : Neo4j
- API publique : REST + OpenAPI
- Interne : monolithe modulaire au départ
- Open Banking : couche d'abstraction multi-fournisseurs
- IA : couche d'abstraction multi-LLM

## Commandes futures

```bash
docker compose up
make test
make lint
make docs
```

Ce dépôt est une v0 documentaire et structurelle. Le code applicatif sera ajouté progressivement.
