# Import CSV bancaire

## Format attendu v3

Le CSV doit contenir les colonnes suivantes :

```csv
date,label,amount,currency
2026-06-25,CARREFOUR MARKET,-82.31,EUR
```

## Règles de validation

- `date` obligatoire au format `YYYY-MM-DD`.
- `label` obligatoire.
- `amount` obligatoire et numérique.
- `currency` obligatoire, 3 caractères recommandés.

## Réponse API

Le système retourne :

- nombre de lignes détectées ;
- nombre de lignes valides ;
- nombre de lignes invalides ;
- détail ligne par ligne.

## Limite v3

Les transactions ne sont pas encore écrites en base.
