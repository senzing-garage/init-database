# init-database examples

## Command line examples

### Command line example using command line option

1. :pencil2: SQLite database.
   Example:

    ```console
    rm -rf /tmp/sqlite
    touch /tmp/sqlite/G2C.db
    senzing-tools init-database --database-url sqlite3://na:na@/tmp/sqlite/G2C.db

    ```

1. :pencil2: PostgreSql database.
   Example:

    ```console
    senzing-tools init-database --database-url postgresql://postgres:postgres@postgres.example.com:5432/G2

    ```

### Command line example using environment variable

1. :pencil2: SQLite database.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db
    senzing-tools init-database

    ```

1. :pencil2: PostgreSql database.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@postgres.example.com:5432/G2
    senzing-tools init-database

    ```

### Command line example adding datasources

In these examples, datasources are added to the initial Senzing configuration.

1. :pencil2: Specify datasources to create using single `--datasources` parameter.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools init-database \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER,REFERENCE,WATCHLIST

    ```

1. :pencil2: Specify datasources to create using multiple `--datasources` parameter.
   Example:

    ```console
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools init-database \
        --database-url postgresql://username:password@postgres.example.com:5432/G2 \
        --datasources CUSTOMER \
        --datasources REFERENCE \
        --datasources WATCHLIST

    ```

1. :pencil2: Specify datasources to create using
   [SENZING_TOOLS_DATASOURCES](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)
   environment variable.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2
    export SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST"
    export LD_LIBRARY_PATH=/opt/senzing/er/lib/
    senzing-tools init-database

    ```

## Docker examples

### Docker example using SENZING_TOOLS_DATABASE_URL

#### Initialize SQLite database

1. Make a directory for Docker volume.
   Example:

    ```console
    mkdir /tmp/my-sqlite

    ```

1. Run `senzing/senzing-tools` Docker image with `init-database` command.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=sqlite3://na:na@/tmp/sqlite/G2C.db \
        --rm \
        --volume /tmp/my-sqlite:/tmp/sqlite \
        senzing/senzing-tools init-database

    ```

#### Initialize PostgreSql database

1. **Optional:** If an existing PostgreSql database doesn't exist,
   [bring up a new PostgreSql database in a docker-compose formation](#docker-compose-stack-with-postgresql-database).

1. :pencil2: Identify database URL, if using existing PostgreSql.
   Example:

    ```console
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@postgres.example.com:5432/G2

    ```

1. Identify database URL, if using docker-compose stack.
   *Note:*  Assignment of `LOCAL_IP_ADDRESS` may not work in all cases.
   Example:

    ```console
    export LOCAL_IP_ADDRESS=$(curl --silent https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/gists/find-local-ip-address/find-local-ip-address.py | python3 -)
    export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@${LOCAL_IP_ADDRESS}:5432/er/?sslmode=disable

    ```

1. Run `senzing/senzing-tools` Docker image with `init-database` command.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL \
        --rm \
        senzing/senzing-tools init-database

    ```

### Docker example using SENZING_TOOLS_ENGINE_CONFIGURATION_JSON

Using `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON` supercedes use of `SENZING_TOOLS_DATABASE_URL`.
It can be used to specify multiple databases or non-system locations of senzing binaries.
For more information, see
[SENZING_TOOLS_ENGINE_CONFIGURATION_JSON](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_json).

1. :pencil2: Set `SENZING_TOOLS_ENGINE_CONFIGURATION_JSON`.
    Example:

    ```console
    export SENZING_TOOLS_ENGINE_CONFIGURATION_JSON='{
        "PIPELINE": {
            "CONFIGPATH": "/etc/opt/senzing",
            "RESOURCEPATH": "/opt/senzing/er/resources",
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
        senzing/senzing-tools init-database
    ```

### Docker example adding datasources

Datasources can be added to the initial Senzing configuration.

1. :pencil2: Specify datasources to create using
   [SENZING_TOOLS_DATASOURCES](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_datasources)
   environment variable.
   Example:

    ```console
    docker run \
        --env SENZING_TOOLS_DATABASE_URL=postgresql://username:password@postgres.example.com:5432/G2 \
        --env SENZING_TOOLS_DATASOURCES="CUSTOMER REFERENCE WATCHLIST" \
        senzing/senzing-tools init-database
    ```

## Miscellaneous

### Docker-compose stack with uninitialized PostgreSql database

1. Identify a directory to place docker-compose artifacts.
   The directory specified will be deleted and re-created.
   Example:

    ```console
    export SENZING_DEMO_DIR=~/my-senzing-demo

    ```

1. Bring up the docker-compose stack.
   Example:

    ```console
    export PGADMIN_DIR=${SENZING_DEMO_DIR}/pgadmin
    export POSTGRES_DIR=${SENZING_DEMO_DIR}/postgres
    export RABBITMQ_DIR=${SENZING_DEMO_DIR}/rabbitmq
    export SENZING_VAR_DIR=${SENZING_DEMO_DIR}/var
    export SENZING_UID=$(id -u)
    export SENZING_GID=$(id -g)

    rm -rf ${SENZING_DEMO_DIR:-/tmp/nowhere/for/safety}
    mkdir ${SENZING_DEMO_DIR}
    mkdir -p ${PGADMIN_DIR} ${POSTGRES_DIR} ${RABBITMQ_DIR} ${SENZING_VAR_DIR}
    chmod -R 777 ${SENZING_DEMO_DIR}

    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-versions-stable.sh \
        https://raw.githubusercontent.com/senzing-garage/knowledge-base/main/lists/docker-versions-stable.sh
    source ${SENZING_DEMO_DIR}/docker-versions-stable.sh
    curl -X GET \
        --output ${SENZING_DEMO_DIR}/docker-compose.yaml \
        "https://raw.githubusercontent.com/senzing-garage/docker-compose-demo/main/resources/postgresql/docker-compose-postgresql-uninitialized.yaml"

    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose up

    ```

1. Cleanup.
   Example:

    ```console
    cd ${SENZING_DEMO_DIR}
    sudo --preserve-env docker-compose down

    ```

### View PostgreSQL database

`pgadmin` is deployed with
[Docker-compose stack with uninitialized PostgreSql database](#docker-compose-stack-with-uninitialized-postgresql-database)

1. Visit [localhost:9171](http://localhost:9171).
   For the initial login, review the instructions at the top of the web page.
   For server password information, see the `POSTGRESQL_POSTGRES_PASSWORD` value in `${SENZING_DEMO_DIR}/docker-compose.yaml`.
   Usually, it's "postgres".

### View SQLite database

The `coleifer/sqlite-web` Docker container can be used to view a SQLite database.

1. Run Docker container.
   Example:

    ```console
    docker run \
        --env SQLITE_DATABASE=G2C.db \
        --interactive \
        --publish 9174:8080 \
        --rm \
        --tty \
        --volume /tmp/sqlite:/data \
        coleifer/sqlite-web

    ```

1. Visit <http://localhost:9174>
