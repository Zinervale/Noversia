# Architecture globale

## Choix initial
Monolithe modulaire pour livrer vite, garder une forte cohérence métier et limiter les coûts.

```text
Clients -> API Gateway -> Core Platform (Go) -> PostgreSQL / Redis / Neo4j -> AI Service (Python) -> LLM Gateway
```

## Modules Core
Identity, Accounts, Transactions, Merchants, Categories, Budget, Goals, Assets, Loans, Recommendations, Notifications, Billing, Audit, Settings.
