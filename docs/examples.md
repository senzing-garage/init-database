# initdatabase examples

## Command line examples

### Command line example - Add datasources

In these examples, datasources are added to the initial Senzing configuration.

1. :pencil2: Specify datasources to create using single `--datasources` parameter.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER,REFERENCE,WATCHLIST
    ```

1. :pencil2: Specify datasources to create using multiple `--datasources` parameter.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER \
        --datasources REFERENCE \
        --datasources WATCHLIST
    ```

1. :pencil2: Specify datasources to create using
   [SENZING_TOOLS_DATASOURCES](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)
   environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST"
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase
    ```

## Docker examples

### Docker example - Using SENZING_TOOLS_ENGINE_CONFIGURATION_JSON

Using `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` superceeds use of `SENZING_TOOLS_DATABASE_URL`.
It can be used to specify multiple databases or non-system locations of senzing binaries.
For more information, see
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json).

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

1. Run `senzing/senzing-tools` Docker container.
    Example:

    ```console
    docker run \
        --env SENZING_TOOLS_ENGINE_CONFIGURATION_JSON \
        --rm \
        senzing/senzing-tools initdatabase
    ```

### Docker example - Add datasources

Datasources can be added to the initial Senzing configuration.

1. :pencil2: Specify datasources to create using
   [SENZING_TOOLS_DATASOURCES](https://github.com/Senzing/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)
   environment variable.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST" \
        senzing/senzing-tools initdatabase
    ```
