# initdatabase examples

## Comm

### XXX

1. :pencil2: Specifying datasources to create.
   Examples:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER,REFERENCE,WATCHLIST
    ```

    or

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    senzing-tools initdatabase \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER \
        --datasources REFERENCE \
        --datasources WATCHLIST
    ```

### xxxx

1. :pencil2: Specifying datasources to create.
   Examples:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST"
    export LD_LIBRARY_PATH=/opt/senzing/g2/lib/
    initdatabase
    ```

### qq

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

### zzzz

1. :pencil2: Run `senzing/initdatabase` specifying datasources to add.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST" \
        senzing/initdatabase
    ```
