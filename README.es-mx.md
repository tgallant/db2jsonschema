# db2jsonschema

![prueba](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main) ![golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Una utilidad para generar definiciones de esquema JSON desde tablas de base de datos.

## Instalación

Ir a la [página de lanzamientos](https://github.com/tgallant/db2jsonschema/releases) para descargar binarios para su sistema.

Compruebe el [DockerHub página](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) para hacer uso de `db2jsonschema` como imagen de docker.

## Usar

`db2jsonschema` se puede utilizar como una aplicación de línea de comandos independiente o importado como una biblioteca en otro paquete de golang.

### Línea de comandos

Por defecto, el comando `db2jsonschema` devolverá un solo documento que contenga todas las definiciones de tabla.

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

Puede guardar el esquema en un archivo mediante la tubería estándar hacia fuera a la deseada ubicación.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db > birds.json
```

El formato de salida se puede cambiar a yaml pasando la opción `--format`.

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

En lugar de emitir un documento único a la norma también es posible escribe cada archivo de esquema a un directorio pasando la opción `--outdir`.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --outdir ./schemas
```

Esto también se puede combinar con la opción `--formate.`

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --format yaml \
  --outdir ./schemas
```

También hay `--include` y `--excluir` los indicadores disponibles para controlar que se utilizarán tablas para la generación de esquemas.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Para pasar en múltiples `--include` o `--excluir` valores proporcionan una coma lista separada o utilice múltiples banderas, uno para cada valor.

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

Las opciones `--schematype` y `--idtemplate` están disponibles para modificar el `$esquema` y `$id` valores añadidos a cada esquema.

`--idtemplate` es una cadena de plantilla que pasa en valores para Nombre y Formato. Estos valores pueden ser utilizados o ignorados según sea necesario.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --schematype https://json-schema.org/draft/2020-12/schema \
  --idtemplate https://example.com/schemas/{{ .Name }}.{{ .Format }} \
  --format yaml \
  --outdir ./schemas
```

### Biblioteca

Aquí hay un ejemplo de importación de `db2jsonschema` como biblioteca y su base El uso de la

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

## Conductores

Hay controladores para conectarse a diferentes backends de base de datos.

- SQLite
- MySQL
- PostgreSQL(WIP)

## Contribuyendo

### Construir

Para crear un binario haga lo siguiente.

```bash
make build
```

Para crear y ejecutar un comando haga lo siguiente.

```bash
./run.sh --driver sqlite3 --dburl ./test.db --format yaml
```

### Prueba

Para ejecutar las pruebas haga lo siguiente.

```bash
make test
```

Para ejecutar todos los cheques haga lo siguiente.

```bash
make test_all
```

Para ejecutar todas las comprobaciones dentro de un entorno de linux containerized haga el Siguiente.

```bash
make ci
```

Para ejecutar todas las acciones de github localmente hacer lo siguiente (requiere https://github.com/nektos/act).

```bash
make actions
```

Para iniciar un servidor MySQL local, ejecute lo siguiente.

```bash
make mysql
```

Para hacer una nueva versión realice lo siguiente (requiere https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```
