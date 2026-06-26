# Backend Go

## Structure v5

```text
cmd/api/main.go
internal/modules/transactions/
  handler.go
  service.go
  repository.go
  parser.go
  models.go
  parser_test.go
```

## Principe

`main.go` câble les dépendances.

Les modules portent la logique métier.
