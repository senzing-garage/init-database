# initdatabase

## :warning: WARNING: initdatabase is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

The Senzing `initdatabase` is used
to initialize databases with a Senzing schema and a default Senzing configuration.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/initdatabase.svg)](https://pkg.go.dev/github.com/senzing/initdatabase)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/initdatabase)](https://goreportcard.com/report/github.com/senzing/initdatabase)
[![go-test.yaml](https://github.com/Senzing/initdatabase/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/initdatabase/actions/workflows/go-test.yaml)

## Overview

Senzing `initdatabase` performs the following:

1. Creates a Senzing database schema from a file of SQL statements found in `/opt/senzing/g2/resources/schema`.
   The file chosen depends on the database engine specified in the protocol section of `SENZING_TOOLS_DATABASE_URL`
   or the database(s) specified in `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
1. Creates a Senzing configuration in the database based on the contents of `/opt/senzing/g2/resources/templates/g2config.json`
1. *Optionally:* Adds datasources to the Senzing configuration.

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
initdatabase [flags]
```

For options and flags, see
[hub.senzing.com/initdatabase](https://hub.senzing.com/initdatabase/) or run:

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
initdatabase --help
```

### Using command line options

1. :pencil2: Specifying database.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase --database-url postgresql://username:password@postgres.example.com:5432/G2
    ```

1. :pencil2: Specifying datasources to create.
   Examples:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER,REFERENCE,WATCHLIST
    ```

    or

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER \
        --datasources REFERENCE \
        --datasources WATCHLIST
    ```

### Using environment variables

1. :pencil2: Specifying database.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase
    ```

1. :pencil2: Specifying datasources to create.
   Examples:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST"
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase
    ```

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. :pencil2: Run `senzing/initdatabase`.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        senzing/initdatabase
    ```

1. *Alternative:* Using `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` environment variable.

    1. :pencil2: Set `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
       Example:

        ```console
        export SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='{
            "PIPELINE": {
                "CONFIGPATH": "/etc/opt/senzing",
                "RESOURCEPATH": "/opt/senzing/g2/resources",
                "SUPPORTPATH": "/opt/senzing/data"
            },
            "SQL": {
                "CONNECTION": "postgresql://username:password@postgres.example.com:5432:G2"
            }
        }'
        ```

    1. Run `senzing/initdatabase`.
       **Note:** `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` superceeds use of `SENZING_TOOLS_DATABASE_URL`.
       Example:

        ```console
        docker run --env SENZING_TOOLS_ENGINE_CONFIGURATION_JSON senzing/initdatabase
        ```

1. :pencil2: Run `senzing/initdatabase` specifying datasources to add.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST" \
        senzing/initdatabase
    ```

### Parameters

- **[SENZING_TOOLS_DATABASE_URL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url)**
- **[SENZING_TOOLS_DATASOURCES](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level)**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name)**
- **[SENZING_TOOLS_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level)**

## Error prefixes

Error identifiers are in the format `senzing-PPPPnnnn` where:

`P` is a prefix used to identify the package.
`n` is a location within the package.

Prefixes:

1. `6501` - initializer
1. `6502` - senzingconfig
1. `6503` - senzingschema

## References

- [Development](docs/development.md)
