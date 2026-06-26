# Modèle de données v8

## Merchants

La table `merchants` stocke :

- nom affiché ;
- nom normalisé.

## Transactions

Les transactions peuvent maintenant pointer vers un marchand via `merchant_id`.

## Suggestions

Les suggestions de règles restent calculées dynamiquement à partir des corrections manuelles.
