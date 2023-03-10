# initdatabase

## :warning: WARNING: initdatabase is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`initdatabase` is a command in the
[senzing-tools](https://github.com/Senzing/senzing-tools)
suite of tools.
This command initialize databases with a Senzing schema and a default Senzing configuration.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/initdatabase.svg)](https://pkg.go.dev/github.com/senzing/initdatabase)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/initdatabase)](https://goreportcard.com/report/github.com/senzing/initdatabase)
[![go-test.yaml](https://github.com/Senzing/initdatabase/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/initdatabase/actions/workflows/go-test.yaml)

## Overview

`initdatabase` performs the following:

1. Creates a Senzing database schema from a file of SQL statements found in `/opt/senzing/g2/resources/schema`.
   The file chosen depends on the database engine specified in the protocol section of `SENZING_TOOLS_DATABASE_URL`
   or the database(s) specified in `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
1. Creates a Senzing configuration in the database based on the contents of `/opt/senzing/g2/resources/templates/g2config.json`
1. *Optionally:* Adds datasources to the Senzing configuration.

## Install

1. The `servegrpc` command is installed with the
   [senzing-tools](https://github.com/Senzing/senzing-tools)
   suite of tools.
   See senzing-tools [install](https://github.com/Senzing/senzing-tools#install).

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
senzing-tools initdatabase [flags]
```

1. For options and flags:
    1. [Online documentation](https://hub.senzing.com/senzing-tools/senzing-tools_initdatabase.html)
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        senzing-tools initdatabase --help
        ```

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase --database-url postgresql://username:password@postgres.example.com:5432/G2
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. :pencil2: Run `senzing/initdatabase`.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --rm \
        senzing/senzing-tools initdatabase
    ```

1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_DATABASE_URL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url)**
- **[SENZING_TOOLS_DATASOURCES](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level)**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name)**
- **[SENZING_TOOLS_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level)**

## References

- [Command reference](https://hub.senzing.com/senzing-tools/senzing-tools_initdatabase.html)
- [Development](docs/development.md)
- [Errors](docs/errors.md)
- [Examples](docs/examples.md)
