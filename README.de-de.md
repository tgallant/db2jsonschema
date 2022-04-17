# db2jsonschema

![test](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main) ![golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Ein Dienstprogramm zur Generierung von JSON Schema aus Datenbanktabellen.

## Installation

Zurück zur [Seite](https://github.com/tgallant/db2jsonschema/releases) um Binaries für Ihr System herunterzuladen.

Prüfen Sie den [DockerHub Seite](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) zu machen Verwendung von `db2jsonschema` als docker

## Nutzung

`db2jsonschema` kann als eigenständige Kommandozeilenanwendung verwendet oder importiert werden als Bibliothek in einem anderen golang

### Befehlszeile

Standardmäßig gibt der Befehl `db2jsonschema` ein einzelnes Dokument zurück, das alle Tabellendefinitionen.

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

Sie können das Schema in einer Datei speichern, indem Sie Standard auf die gewünschte Datei piping Lage

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db > birds.json
```

Das Ausgabeformat kann durch Übergeben der Option `--format` in yaml geändert werden.

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

Anstatt ein einzelnes Dokument auf Standardausgabe auszugeben, ist es auch möglich, schreiben Sie jede Schemadatei in ein Verzeichnis durch Übergeben der Option `--outdir`.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --outdir ./schemas
```

Dies kann auch mit der Option `--format` kombiniert werden.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --format yaml \
  --outdir ./schemas
```

Es gibt auch `--inclusive` und `---ausschließen` Flags zur Steuerung der Tabellen werden für die Schemagenerierung verwendet.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Um in mehreren `--Include-``` oder -ausschließen Werte zu übergeben, geben Sie entweder ein Komma getrennte Liste oder verwenden Sie mehrere Flags, eines für jeden Wert.

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

Die Optionen `--schematype` und `--idtemplate` stehen zur Verfügung, um die `$schema` und `$id` Werte zu jedem Schema hinzugefügt.

`--idtemplate` ist eine template die in Werten für Name und Format übergeht. Diese Werte können bei Bedarf verwendet oder ignoriert werden.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --schematype https://json-schema.org/draft/2020-12/schema \
  --idtemplate https://example.com/schemas/{{ .Name }}.{{ .Format }} \
  --format yaml \
  --outdir ./schemas
```

### Bibliothek

Hier ist ein Beispiel für den Import `db2jsonschema` als Bibliothek und seine grundlegende Nutzung.

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

## Fahrer

Es gibt Treiber für die Verbindung mit verschiedenen Datenbank-Backends.

- SQLite
- MySQL
- PostgreSQL (WIP)

## Mitwirkend

### Bauen

Um eine Binärdatei zu erstellen, führen Sie Folgendes aus.

```bash
make build
```

Um einen Befehl zu erstellen und auszuführen, führen Sie Folgendes aus.

```bash
./run.sh --driver sqlite3 --dburl ./test.db --format yaml
```

### Testen

Um die Tests durchzuführen, führen Sie Folgendes aus.

```bash
make test
```

Um alle Kontrollen auszuführen, führen Sie Folgendes aus.

```bash
make test_all
```

Um alle Kontrollen innerhalb einer containerisierten linux-Umgebung durchzuführen, führen Sie die - und nachfolgend.

```bash
make ci
```

Um alle github lokal auszuführen führen führen (erfordert https://github.com/nektos/act).

```bash
make actions
```

Um einen lokalen MySQL-Server zu starten, führen Sie die folgenden aus.

```bash
make mysql
```

Um eine neue Version durchzuführen, führen Sie Folgendes (erfordert https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```