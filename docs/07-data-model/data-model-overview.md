# Modèle de données v7

## Nouvelle table

### transaction_enrichments

Historise les changements apportés à une transaction.

Colonnes :
- transaction_id
- enrichment_type
- previous_value
- new_value
- source
- reason
- created_at

## Objectif

Garantir la traçabilité des corrections et préparer l'apprentissage.
