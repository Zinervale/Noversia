# Marchands

## Objectif

Regrouper les transactions par commerçant réel.

## Approche v8

Le marchand est extrait du libellé bancaire.

Exemples :

| Libellé | Marchand |
|---|---|
| CARREFOUR MARKET 1234 | CARREFOUR |
| NETFLIX.COM | NETFLIX |
| TOTAL ENERGIES | TOTAL |

## Limite

La détection reste simple. Les futures versions intégreront :
- alias marchands ;
- apprentissage utilisateur ;
- enrichissement IA ;
- regroupement multi-libellés.
