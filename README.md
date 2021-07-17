# db2jsonschema

A utility for generating JSON Schema definitions from database tables.

## Usage

`db2jsonschema` can be used as a standalone command line application or imported
as a library in another golang package.

### Command Line

By default the `db2jsonschema` command will return a single document containing
all of the table definitions.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://example.com/birds.schema.json",
  "title": "Birds",
  "definitions": {
    "Example": {
      "id": {
        "name": "id",
        "type": "number"
      },
      "name": {
        "name": "name",
        "type": "string"
      }
    }
  }
}
```

You can save the schema to a file by piping standard out to the desired
location.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db > birds.json
```

The output format can be changed to yaml by passing the `--format` option.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --format yaml
$schema: https://json-schema.org/draft/2020-12/schema
$id: https://example.com/birds.schema.json
title: Birds
type: object
properties:
  id:
    name: id
    type: number
  name:
    name: name
    type: string
```

Instead of outputting a single document to standard out it is also possible to
write each schema file to a directory by passing the `--outdir` option.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --outdir ./schemas
```

This can also be combined with the `--format` option.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --format yaml \
  --outdir ./schemas
```

There are also `--include` and `--exclude` flags available to control which
tables will be used for schema generation.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include Birds
  --format yaml \
  --outdir ./schemas
```

To pass in multiple `--include` or `--exclude` values either provide a comma
separated list or use multiple flags, one for each value.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --exclude Locations,BirdWatchers
  --format yaml \
  --outdir ./schemas
```

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --exclude Locations
  --exclude BirdWatchers
  --format yaml \
  --outdir ./schemas
```

### Library

Here is an example of importing `db2jsonschema` as a library and its basic
usage.

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

## Drivers

There are drivers for connecting to different database backends.

- SQLite
- MySQL (WIP)
- PostgreSQL(WIP)

## Contributing

### Build

To build a binary do the following.

```bash
go build -o db2jsonschema cmd/db2jsonschema/main.go
```

To build and run a command do the following.

```bash
go run cmd/db2jsonschema/main.go --driver sqlite3 --dburl ./test.db --format yaml
```

### Test

To run the tests do the following.

```bash
go test -v ./...
```
