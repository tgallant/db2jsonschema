# db2jsonschema

![prueba](https://github.com/tgallant/db2jsonschema/actions/workflows/test.yaml/badge.svg?branch=main) ![golangci](https://github.com/tgallant/db2jsonschema/actions/workflows/lint.yaml/badge.svg?branch=main) ![shellcheck](https://github.com/tgallant/db2jsonschema/actions/workflows/shellcheck.yaml/badge.svg?branch=main)

Una utilidad para generar definiciones de esquemas JSON a partir de tablas de bases de datos.

## Instalación

[Dirígete a la página de lanzamientos](https://github.com/tgallant/db2jsonschema/releases)
 para descargar binarios para su sistema.

Compruebe el [DockerHub
](https://hub.docker.com/repository/docker/tgallant/db2jsonschema) página para hacer
 El uso de `db2jsonschema` como una imagen docker.

## Uso

`db2jsonschema` puede ser utilizado como una aplicación independiente de línea de comandos o importado
 como una biblioteca en otro paquete de golang.

### Línea de comandos

Por defecto, el comando `db2jsonschema` devolverá un único document que contenga
 todas las definiciones de la tabla.

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

Puede guardar el esquema en un archivo canalizando el estándar al formato deseado
 Localización

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

En lugar de producir un solo documento para estandarizar, también es posible
 Escribe cada archivo de esquema en un directorio pasando la opción `--outdir`.

```bash
db2jsonschema --driver sqlite3 --dburl ./exotic_birds.db --outdir ./schemas
```

Esto también se puede combinar con la opción `--format`.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --format yaml \
  --outdir ./schemas
```

También hay indicadores `--include` y `--exclude` disponibles para controlar qué
 Las tablas se utilizarán para la generación de esquemas.

```bash
db2jsonschema \
  --driver sqlite3 \
  --dburl ./exotic_birds.db \
  --include birds
  --format yaml \
  --outdir ./schemas
```

Para pasar en múltiples valores `--include` o `--exclude` proporcione una coma
 lista separada o utilizar múltiples banderas, una para cada valor.

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

Las opciones `--schematype` y `--idtemplate` están disponibles para modificar el
 Los valores `$schema` y `$id` se agregan a cada esquema.

`--idtemplate` es una cadena de plantilla que pasa los valores de Nombre y Formato.
 Estos valores pueden ser usados o ignorados según sea necesario.

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

Aquí hay un ejemplo de importación de `db2jsonschema` como biblioteca y su
 uso.

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

Hay controladores para conectarse a diferentes backends de bases de datos.

- SQLite
- MySQL
- PostgreSQL (WIP)

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

Para ejecutar todas las comprobaciones, haga lo siguiente.

```bash
make test_all
```

Para ejecutar todas las comprobaciones dentro de un entorno linux en contenedores, haga
 Siguiente.

```bash
make ci
```

Para ejecutar todas las acciones de github localmente, haga lo siguiente (requiere
 https://github.com/nektos/act).

```bash
make actions
```

Para iniciar un servidor local MySQL ejecute lo siguiente.

```bash
make mysql
```

Para hacer una nueva versión realice lo siguiente (requiere
 https://docs.npmjs.com/cli/v6/using-npm/semver).

```bash
# possible options for SEMVER are:
# major, minor, patch, premajor, preminor, prepatch, or prerelease
make release SEMVER=patch
```
