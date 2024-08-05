# init-database

If you are beginning your journey with [Senzing],
please start with [Senzing Quick Start guides].

You are in the [Senzing Garage] where projects are "tinkered" on.
Although this GitHub repository may help you understand an approach to using Senzing,
it's not considered to be "production ready" and is not considered to be part of the Senzing product.
Heck, it may not even be appropriate for your application of Senzing!

## :warning: WARNING: init-database is still in development :warning: _

At the moment, this is "work-in-progress" with Semantic Versions of `0.n.x`.
Although it can be reviewed and commented on,
the recommendation is not to use it yet.

## Synopsis

`init-database` is a command in the
[senzing-tools]
suite of tools.
This command initialize databases with a Senzing schema and a default Senzing configuration.

[![Go Reference Badge]][Package reference]
[![Go Report Card Badge]][Go Report Card]
[![License Badge]][License]
[![go-test-linux.yaml Badge]][go-test-linux.yaml]
[![go-test-darwin.yaml Badge]][go-test-darwin.yaml]
[![go-test-windows.yaml Badge]][go-test-windows.yaml]

[![golangci-lint.yaml Badge]][golangci-lint.yaml]

## Overview

`init-database` performs the following:

1. Creates a Senzing database schema from a file of SQL statements.
   The SQL file is identified by the [SENZING_TOOLS_SQL_FILE] parameter.
   The default file name depends on the database engine specified in the
   protocol section of [SENZING_TOOLS_DATABASE_URL]
   or the database(s) specified in [SENZING_TOOLS_ENGINE_CONFIGURATION_JSON].
   The default file location depends on the Senzing engine configuration JSON's `PIPELINE`.`RESOURCEPATH` value.
1. Creates a Senzing configuration in the database based on the contents
   of the file specified by the [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE] parameter.
   The default file location is based on the Senzing engine configuration JSON's `PIPELINE`.`RESOURCEPATH` value.
1. *Optionally:* Adds datasources to the initial Senzing configuration via the [SENZING_TOOLS_DATASOURCES] parameter.

## Install

1. The `init-database` command is installed with the [senzing-tools] suite of tools.
   See senzing-tools [install].

## Use

```console
export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
senzing-tools init-database [flags]
```

1. For options and flags:
    1. [Online documentation]
    1. Runtime documentation:

        ```console
        export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
        senzing-tools init-database --help
        ```

1. In addition to the following simple usage examples, there are additional [Examples].

### Using command line options

1. :pencil2: Specify database using command line option.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools init-database \
        --database-url postgresql://username:password@postgres.example.com:5432/G2
    ```

1. See [Parameters] for additional parameters.

### Using environment variables

1. :pencil2: Specify database using environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools init-database
    ```

1. See [Parameters] for additional parameters.

### Using Docker

This usage shows how to initialze a database with a Docker container.

1. :pencil2: Run `senzing/senzing-tools`.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_COMMAND=init-database \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --rm \
        senzing/senzing-tools
    ```

1. See [Parameters] for additional parameters.

### Parameters

- **[SENZING_TOOLS_COMMAND]**
- **[SENZING_TOOLS_CONFIG_PATH]**
- **[SENZING_TOOLS_CONFIGURATION]**
- **[SENZING_TOOLS_DATABASE_URL]**
- **[SENZING_TOOLS_DATASOURCES]**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]**
- **[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]**
- **[SENZING_TOOLS_ENGINE_LOG_LEVEL]**
- **[SENZING_TOOLS_ENGINE_MODULE_NAME]**
- **[SENZING_TOOLS_LICENSE_STRING_BASE64]**
- **[SENZING_TOOLS_LOG_LEVEL]**
- **[SENZING_TOOLS_OBSERVER_ORIGIN]**
- **[SENZING_TOOLS_OBSERVER_URL]**
- **[SENZING_TOOLS_RESOURCE_PATH]**
- **[SENZING_TOOLS_SENZING_DIRECTORY]**
- **[SENZING_TOOLS_SQL_FILE]**
- **[SENZING_TOOLS_SUPPORT_PATH]**

## References

- [Command reference]
- [Development]
- [Errors]
- [Examples]

[Development]: docs/development.md
[Errors]: docs/errors.md
[Examples]: docs/examples.md
[Go Reference Badge]: https://pkg.go.dev/badge/github.com/senzing-garage/template-go.svg
[Go Report Card Badge]: https://goreportcard.com/badge/github.com/senzing-garage/template-go
[Go Report Card]: https://goreportcard.com/report/github.com/senzing-garage/template-go
[License Badge]: https://img.shields.io/badge/License-Apache2-brightgreen.svg
[License]: https://github.com/senzing-garage/template-go/blob/main/LICENSE
[Online documentation]: https://hub.senzing.com/senzing-tools/senzing-tools_init-database.html
[Package reference]: https://pkg.go.dev/github.com/senzing-garage/template-go
[Parameters]: #parameters
[SENZING_TOOLS_COMMAND]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_command
[SENZING_TOOLS_CONFIGURATION]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_configuration
[SENZING_TOOLS_CONFIG_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_config_path
[SENZING_TOOLS_DATABASE_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_database_url
[SENZING_TOOLS_DATASOURCES]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources
[SENZING_TOOLS_ENGINE_CONFIGURATION_FILE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_file
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json
[SENZING_TOOLS_ENGINE_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_log_level
[SENZING_TOOLS_ENGINE_MODULE_NAME]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_module_name
[SENZING_TOOLS_LICENSE_STRING_BASE64]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_license_string_base64
[SENZING_TOOLS_LOG_LEVEL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_log_level
[SENZING_TOOLS_OBSERVER_ORIGIN]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_observer_origin
[SENZING_TOOLS_OBSERVER_URL]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_observer_url
[SENZING_TOOLS_RESOURCE_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_resource_path
[SENZING_TOOLS_SENZING_DIRECTORY]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_senzing_directory
[SENZING_TOOLS_SQL_FILE]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_sql_file
[SENZING_TOOLS_SUPPORT_PATH]: https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_support_path
[Senzing Garage]: https://github.com/senzing-garage
[Senzing Quick Start guides]: https://docs.senzing.com/quickstart/
[Senzing]: https://senzing.com/
[go-test-darwin.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml/badge.svg
[go-test-darwin.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-darwin.yaml
[go-test-linux.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml/badge.svg
[go-test-linux.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-linux.yaml
[go-test-windows.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml/badge.svg
[go-test-windows.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/go-test-windows.yaml
[golangci-lint.yaml Badge]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml/badge.svg
[golangci-lint.yaml]: https://github.com/senzing-garage/template-go/actions/workflows/golangci-lint.yaml
[install]: https://github.com/senzing-garage/senzing-tools#install
[senzing-tools]: https://github.com/senzing-garage/senzing-tools
