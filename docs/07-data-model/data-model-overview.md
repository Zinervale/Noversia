# Modèle de données v4

## Nouvelles tables

### import_batches

Représente un fichier importé.

### import_rows

Représente chaque ligne du fichier avec son statut.

### transactions.source_hash

Empreinte fonctionnelle pour limiter les doublons.

## Principe

Une ligne invalide est conservée dans `import_rows`, mais ne crée pas de transaction.
