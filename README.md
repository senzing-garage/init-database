# init-database

## :warning: WARNING: init-database is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`init-database` is a command in the
[senzing-tools](https://github.com/Senzing/senzing-tools)
suite of tools.
This command initialize databases with a Senzing schema and a default Senzing configuration.

[![Go Reference](https://pkg.go.dev/badge/github.com/senzing/init-database.svg)](https://pkg.go.dev/github.com/senzing/init-database)
[![Go Report Card](https://goreportcard.com/badge/github.com/senzing/init-database)](https://goreportcard.com/report/github.com/senzing/init-database)
[![go-test.yaml](https://github.com/Senzing/init-database/actions/workflows/go-test.yaml/badge.svg)](https://github.com/Senzing/init-database/actions/workflows/go-test.yaml)
[![License](https://img.shields.io/badge/License-Apache2-brightgreen.svg)](https://github.com/Senzing/init-database/blob/main/LICENSE)

## Overview

`init-database` performs the following:

1. Creates a Senzing database schema from a file of SQL statements.
   The SQL file is identified by the `SENZING_TOOLS_SQL_FILE` parameter.
   The default file depends on the database engine specified in the
   protocol section of `SENZING_TOOLS_DATABASE_URL`
   or the database(s) specified in `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
1. Creates a Senzing configuration in the database based on the contents
   of the file specified by the `SENZING_TOOLS_ENGINE_CONFIGURATION_FILE` parameter.
1. *Optionally:* Adds datasources to the initial Senzing configuration
   via the `SENZING_TOOLS_DATASOURCES` parameter.

## Install

1. The `init-database` command is installed with the
   [senzing-tools](https://github.com/Senzing/senzing-tools)
   suite of tools.
   See senzing-tools [install](https://github.com/Senzing/senzing-tools#install).

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
senzing-tools init-database [flags]
```

1. For options and flags:
    1. [Online documentation](https://hub.senzing.com/senzing-tools/senzing-tools_init-database.html)
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        senzing-tools init-database --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples](docs/examples.md).

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools init-database \
        --database-url postgresql://username:password@postgres.example.com:5432/G2
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools init-database
    ```

1. See [Parameters](#parameters) for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. :pencil2: Run `senzing/init-database`.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --rm \
        senzing/senzing-tools init-database
    ```

1. See [Parameters](#parameters) for additional parameters.

### Parameters

- **[SENZING_TOOLS_CONFIG_PATH](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_config_path)**
- **[SENZING_TOOLS_CONFIGURATION](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_configuration)**
- **[SENZING_TOOLS_DATABASE_URL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url)**
- **[SENZING_TOOLS_DATASOURCES](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_FILE](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_file)**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json)**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level)**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name)**
- **[SENZING_TOOLS_LICENSE_STRING_BASE64](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_license_string_base64)**
- **[SENZING_TOOLS_LOG_LEVEL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level)**
- **[SENZING_TOOLS_OBSERVER_ORIGIN](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_observer_origin)**
- **[SENZING_TOOLS_OBSERVER_URL](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_observer_url)**
- **[SENZING_TOOLS_RESOURCE_PATH](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_resource_path)**
- **[SENZING_TOOLS_SENZING_DIRECTORY](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_senzing_directory)**
- **[SENZING_TOOLS_SQL_FILE](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_sql_file)**
- **[SENZING_TOOLS_SUPPORT_PATH](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_support_path)**

## References

- [Command reference](https://hub.senzing.com/senzing-tools/senzing-tools_init-database.html)
- [Development](docs/development.md)
- [Errors](docs/errors.md)
- [Examples](docs/examples.md)
