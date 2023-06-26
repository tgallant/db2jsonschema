# db2jsonschema

![Golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main)![](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![Shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Ein Dienstprogramm zum Generieren von JSON-Schemadefinitionen aus Datenbanktabellen.

## Installation

[Gehe zur Releaseseite](https://github.com/tgallant/db2jsonschema/releases)
 Binärdateien für Ihr System herunterzuladen.

Überprüfen Sie den [DockerHub
](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) Seite zu machen
 `db2jsonschema` als Docker-Image.

## Verwendung

`db2jsonschema` kann als eigenständige Befehlszeilenanwendung verwendet oder importiert werden
 als Bibliothek in einem anderen Golang-Paket.

### Kommandozeile

Standardmäßig gibt der Befehl `db2jsonschema` ein einzelnes Dokument zurück, das
 alle Tabellendefinitionen.

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

Sie können das Schema in einer Datei speichern, indem Sie den Standard an die gewünschte
 Standort

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db > birds.json
```

Das Ausgabeformat kann durch Übergabe der Option `--format` in yaml geändert werden.

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

Anstatt ein einzelnes Dokument auf Standard Out auszugeben, ist es auch möglich,
 schreibt jede Schema-Datei in ein Verzeichnis, indem man die Option `--outdir` übergibt.

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

Es gibt auch `--include` und `--exclude` Flags zur Verfügung, um zu steuern, welche
 Tabellen werden für die Schemagenerierung verwendet.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Um mehrere `--include`- oder `--exclude`-Werte zu übergeben, geben Sie entweder ein Komma an
 Sie können eine separate Liste oder mehrere Flags verwenden, einen für jeden Wert.

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

Die Optionen `--schematype` und `--idtemplate` stehen zur Verfügung, um die
 `$schema` und `$id` Werte werden zu jedem Schema hinzugefügt.

`--idtemplate` ist eine Vorlagenzeichenfolge, die Werte für Name und Format übergibt.
 Diese Werte können bei Bedarf verwendet oder ignoriert werden.

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

Hier ist ein Beispiel für den Import von `db2jsonschema` als Bibliothek und seine
 Nutzung.

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

## Treiber

Es gibt Treiber für die Verbindung zu verschiedenen Datenbank-Backends.

- SQLite
- MySQL
- PostgreSQL (WIP)

## Mitmachen

### Bauen

Um eine Binärdatei zu erstellen, gehen Sie wie folgt vor.

```bash
make build
```

Um einen Befehl zu erstellen und auszuführen, gehen Sie folgendermaßen vor:

```bash
./run.sh --driver sqlite3 --dburl ./test.db --format yaml
```

### Testen

Um die Tests auszuführen, gehen Sie wie folgt vor.

```bash
make test
```

Um alle Prüfungen auszuführen, gehen Sie wie folgt vor.

```bash
make test_all
```

Um alle Prüfungen in einer containerisierten Linux-Umgebung auszuführen,
 folgen.

```bash
make ci
```

Um alle github-Aktionen lokal auszuführen, gehen Sie folgendermaßen vor (erfordert
 https://github.com/nektos/act).

```bash
make actions
```

Um einen lokalen MySQL-Server zu starten, führen Sie die folgenden Schritte aus:

```bash
make mysql
```

Um eine neue Version zu erstellen, führen Sie folgende Schritte aus (erfordert
 https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```
