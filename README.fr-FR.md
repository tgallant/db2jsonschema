# db2jsonschema

![test](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main) ![golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Un utilitaire pour générer des définitions de schéma JSON à partir de tables de base de données.

## Installation

Rendez-vous sur la [page des communiqués](https://github.com/tgallant/db2jsonschema/releases)
 pour télécharger des binaires pour votre système.

Vérifiez le [DockerHub
](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) Page à faire
 utilisation de `db2jsonschema` comme image de docker.

## Utilisation

`db2jsonschema` peut être utilisé comme application autonome en ligne de commande ou importé
 comme une bibliothèque dans un autre paquet golang.

### Ligne de commande

Par défaut, la commande `db2jsonschema` renverra un seul document
 toutes les définitions de tableaux.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "definitions.json",
  "title": "Definitions",
  "definitions": {
    "birds": {
      "created_at": {
        "name": "created_at",
        "type": "string",
        "format": "date-time"
      },
      "deleted_at": {
        "name": "deleted_at",
        "type": "string",
        "format": "date-time"
      },
      "id": {
        "name": "id",
        "type": "number"
      },
      "genus": {
        "name": "genus",
        "type": "string"
      },
      "species": {
        "name": "species",
        "type": "string"
      },
      "updated_at": {
        "name": "updated_at",
        "type": "string",
        "format": "date-time"
      }
    }
  }
}
```

Vous pouvez enregistrer le schéma dans un dossier en passant par le standard vers le dossier souhaité
 Localisation.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db > birds.json
```

Le format de sortie peut être changé en yaml en passant l'option `--format`.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --format yaml
$schema: https://json-schema.org/draft/2020-12/schema
$id: birds.yaml
title: Definitions
definitions:
  birds:
    created_at:
      name: created_at
      type: string
      format: date-time
    deleted_at:
      name: deleted_at
      type: string
      format: date-time
    id:
      name: id
      type: number
    genus:
      name: genus
      type: string
    species:
      name: species
      type: string
    updated_at:
      name: updated_at
      type: string
      format: date-time
```

Au lieu de produire un seul document à normaliser, il est également possible de
 écrire chaque schéma dans un répertoire en passant l'option `--outdir`.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --outdir ./schemas
```

Cela peut également être combiné avec l'option `--format`.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --format yaml \
  --outdir ./schemas
```

Il existe également des drapeaux `--include` et `--exclude` disponibles pour commander qui
 Les tables seront utilisées pour la génération de schémas.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Pour passer dans plusieurs valeurs `--include` ou `--exclude` soit fournir une virgule
 séparée ou utilisez plusieurs drapeaux, un pour chaque valeur.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --exclude locations,bird_watchers \
  --format yaml \
  --outdir ./schemas
```

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --exclude locations \
  --exclude bird_watchers \
  --format yaml \
  --outdir ./schemas
```

Les options `--schematype` et `--idtemplate` sont disponibles pour modifier le
 Les valeurs `$schema` et `$id` ajoutées à chaque schéma.

`--idtemplate` est une chaîne de modèle qui passe dans les valeurs de Name et Format.
 Ces valeurs peuvent être utilisées ou ignorées si nécessaire.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --schematype https://json-schema.org/draft/2020-12/schema \
  --idtemplate https://example.com/schemas/{{ .Name }}.{{ .Format }} \
  --format yaml \
  --outdir ./schemas
```

### Bibliothèque

Voici un exemple d'importation de `db2jsonschema` en tant que bibliothèque et son
 Utilisation.

```golang
package main
import (
    "fmt"
	"github.com/tgallant/db2jsonschema"
)
func main() {
  request := &db2jsonschema.Request{
      Driver:     "sqlite3",
      DataSource: "./test.db",
      Format:     "yaml",
      Outdir:     "./schemas",
  }
  err := request.Perform()
  if err != nil {
      fmt.Println(err)
  }
}
```

## Chauffeurs

Il existe des pilotes pour relier différents backends de base de données.

- SQLite
- MySQL
- PostgreSQL

## Participer

### Créateur

Pour créer un binaire, procédez comme suit.

```bash
make build
```

Pour créer et exécuter une commande, procédez comme suit.

```bash
./run.sh --driver sqlite3 --dburl ./test.db --format yaml
```

### Test

Pour exécuter les tests, procédez comme suit.

```bash
make test
```

Pour exécuter toutes les vérifications, procédez comme suit.

```bash
make test_all
```

Pour exécuter toutes les vérifications à l'intérieur d'un environnement Linux
 Suivant.

```bash
make ci
```

Pour exécuter toutes les actions github localement, procédez comme suit (nécessite
 https://github.com/nektos/act).

```bash
make actions
```

Pour démarrer un serveur MySQL local, exécutez ce qui suit.

```bash
make mysql
```

Pour créer une nouvelle version, procédez comme suit (nécessite
 https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```
