# Backend Go

## v3

Ajout du parsing CSV réel directement dans le handler d'import.

## Prochaine refactorisation

Déplacer le parsing dans un service dédié :

```text
internal/modules/transactions/importer.go
```
