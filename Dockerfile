# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.25.3-bookworm@sha256:51b6b12427dc03451c24f7fc996c43a20e8a8e56f0849dd0db6ff6e9225cc892
ARG IMAGE_FINAL=senzing/senzingsdk-runtime:4.1.0@sha256:e57d751dc0148bb8eeafedb7accf988413f50b54a7e46f25dfe4559d240063e5

ARG SENZING_APT_INSTALL_SETUP_PACKAGE="senzingsdk-setup"

# -----------------------------------------------------------------------------
# Stage: senzingsdk_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS senzingsdk_runtime

# -----------------------------------------------------------------------------
# Stage: builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_BUILDER} AS builder
ENV REFRESHED_AT=2024-08-01
LABEL Name="senzing/go-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.0"

# Run as "root" for system installation.

USER root

# Install packages via apt-get.

RUN apt-get update \
  && apt-get -y --no-install-recommends install \
      libsqlite3-dev \
      wget \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/init-database

# Copy files from prior stage.

COPY --from=senzingsdk_runtime  "/opt/senzing/er/lib/"   "/opt/senzing/er/lib/"
COPY --from=senzingsdk_runtime  "/opt/senzing/er/sdk/c/" "/opt/senzing/er/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/
WORKDIR ${GOPATH}/src/init-database

# Build go program.

RUN make build \
  && go clean -cache -modcache -testcache

# Copy binaries to /output.

RUN mkdir -p /output \
 && cp -R ${GOPATH}/src/init-database/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: senzingsdk
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS senzingsdk
ENV REFRESHED_AT=2024-08-01

ARG SENZING_APT_INSTALL_SETUP_PACKAGE

ENV SENZING_APT_INSTALL_SETUP_PACKAGE=${SENZING_APT_INSTALL_SETUP_PACKAGE}

# Install Senzing package.

RUN apt-get update \
  && apt-get -y --no-install-recommends install ${SENZING_APT_INSTALL_SETUP_PACKAGE} \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

# -----------------------------------------------------------------------------
# Stage: oracle
# -----------------------------------------------------------------------------

# FROM ${IMAGE_FINAL} AS oracle
# ENV REFRESHED_AT=2024-08-01

# RUN apt-get update \
#  && apt-get -y install \
#       curl \
#       unzip

# RUN curl -X GET \
#         --output /tmp/instantclient-basiclite-linux.zip \
#         https://download.oracle.com/otn_software/linux/instantclient/2390000/instantclient-basiclite-linux.x64-23.9.0.25.07.zip

# RUN mkdir /opt/oracle \
#  && cd /opt/oracle \
#  && unzip /tmp/instantclient-basiclite-linux.zip

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2024-08-01
LABEL Name="senzing/init-database" \
      Maintainer="support@senzing.com" \
      Version="0.7.19"
HEALTHCHECK CMD ["/app/healthcheck.sh"]
USER root

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y --no-install-recommends install \
      libsqlite3-dev \
      default-libmysqlclient-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# MySQL support.

RUN wget https://dev.mysql.com/get/Downloads/Connector-ODBC/9.5/mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
 && wget https://dev.mysql.com/get/Downloads/MySQL-9.5/mysql-common_9.5.0-1debian13_amd64.deb \
 && wget https://deb.debian.org/debian/pool/main/m/mysql-8.0/libmysqlclient21_8.0.44-1_amd64.deb \
 && apt-get update \
 && apt-get -y install \
      ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
      ./mysql-common_9.5.0-1debian13_amd64.deb \
      ./libmysqlclient21_8.0.44-1_amd64.deb \
 && rm \
      ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
      ./mysql-common_9.5.0-1debian13_amd64.deb \
      ./libmysqlclient21_8.0.44-1_amd64.deb \
 && rm -rf /var/lib/apt/lists/*

# MS SQL support.

RUN curl -sSL -O https://packages.microsoft.com/config/ubuntu/"$(grep VERSION_ID /etc/os-release | cut -d '"' -f 2)"/packages-microsoft-prod.deb \
 && dpkg -i --force-confnew packages-microsoft-prod.deb /
 && rm packages-microsoft-prod.deb \
 && apt-get update \
 && ACCEPT_EULA=Y apt-get install -y msodbcsql18 \
 && ACCEPT_EULA=Y apt-get install -y mssql-tools18 \
 && apt-get install -y unixodbc-dev

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/init-database /app/init-database
COPY --from=senzingsdk /opt/senzing/er/resources/schema /opt/senzing/er/resources/schema
# COPY --from=oracle /opt/oracle /opt/oracle

# Run as non-root container.

USER 1001

# Runtime environment variables.

# ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/oracle/instantclient_23_9

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
