# db2jsonschema

![test](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main) ![golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Un utilitaire pour générer des définitions JSON Schema à partir de tables de base de données.

## Installation

Rendez-vous sur la [page de sortie](https://github.com/tgallant/db2jsonschema/releases) pour télécharger des binaires pour votre système.

Vérifiez le [DockerHub page](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) à faire utilisation de `db2jsonschema` comme une image docker.

## Utilisation

`db2jsonschema` peut être utilisé comme application de ligne de commande autonome ou importé comme bibliothèque dans un autre paquet golang.

### Ligne de commande

Par défaut, la commande `db2jsonschema` retourne un seul document contenant toutes les définitions du tableau.

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

Vous pouvez enregistrer le schéma dans un fichier en utilisant la norme de piping vers le désiré Emplacement.

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

Au lieu de produire un seul document pour normaliser, il est également possible de écrire chaque fichier schéma dans un répertoire en passant l'option `--outdir`.

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

Il y a aussi `--inclure` et `--exclure` les drapeaux disponibles pour contrôler les tables seront utilisées pour la génération de schéma.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Pour passer en plusieurs valeurs `--inclure` ou `--exclut` soit fournir une virgule liste séparée ou utilisez plusieurs drapeaux, un pour chaque valeur.

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

Les options `--schematype` et `--idtemplate` sont disponibles pour modifier le Les valeurs `$schema` et `$id` sont ajoutées à chaque schéma.

`--idtemplate` est une chaîne de caractères de modèle qui passe dans les valeurs pour Nom et Format. Ces valeurs peuvent être utilisées ou ignorées au besoin.

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

Voici un exemple d'importation `db2jsonschema` comme bibliothèque et de son base l'utilisation.

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

## Les chauffeurs

Il existe des pilotes pour se connecter à différents backend.

- SQLite
- MySQL
- PostgreSQL(WIP)

## Contribuer

### Construire

Pour construire un binaire, procédez comme suit.

```bash
make build
```

Pour construire et exécuter une commande procédez comme suit.

```bash
./run.sh --driver sqlite3 --dburl ./test.db --format yaml
```

### - T'as pas le choix

Pour exécuter les tests, procédez comme suit.

```bash
make test
```

Pour exécuter tous les contrôles, faites les choses suivantes.

```bash
make test_all
```

Pour exécuter toutes les vérifications à l'intérieur d'un environnement linux conteneurisé suivre.

```bash
make ci
```

Pour exécuter toutes les actions github localement faites les opérations suivantes (nécessite https://github.com/nektos/act).

```bash
make actions
```

Pour démarrer un serveur MySQL local, exécutez les opérations suivantes.

```bash
make mysql
```

Pour faire une nouvelle version, exécutez les éléments suivants (nécessite https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```
