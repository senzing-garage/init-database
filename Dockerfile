# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.25.4-bookworm@sha256:7419f544ffe9be4d7cbb5d2d2cef5bd6a77ec81996ae2ba15027656729445cc4
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
ENV REFRESHED_AT=2025-11-15
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
ENV REFRESHED_AT=2025-11-15

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
# ENV REFRESHED_AT=2025-11-15

# RUN apt-get update \
#  && apt-get -y install \
#       curl \
#       unzip

# RUN curl -X GET \
#         --output /tmp/instantclient-basiclite-linux.zip \
#         https://download.oracle.com/otn_software/linux/instantclient/2326000/instantclient-basic-linux.x64-23.26.0.0.0.zip

# RUN mkdir /opt/oracle \
#  && cd /opt/oracle \
#  && unzip /tmp/instantclient-basiclite-linux.zip

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} AS final
ENV REFRESHED_AT=2025-11-15
LABEL Name="senzing/init-database" \
      Maintainer="support@senzing.com" \
      Version="0.7.19"
HEALTHCHECK CMD ["/app/healthcheck.sh"]
USER root

ENV ACCEPT_EULA=y \
    DEBIAN_FRONTEND=noninteractive

# Work-around for apt-get update error.

RUN chmod 1777 /tmp

# Install packages via apt-get.

RUN apt-get update \
 && apt-get -y --no-install-recommends install \
      apt-transport-https \
      curl \
      default-libmysqlclient-dev \
      gnupg \
      libaio-dev \
      libsqlite3-dev \
      lsb-release \
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

# See docker-compose/docker-compose.mssql.yaml's use of mssql-driver-volume to see how to attach MS SQL drivers.

# Tricky code: This is installing the MS SQL driver for Debian 12 instead of 13,
# because MS doesn't have Debian 13 builds yet, as of 2025-11-12.

# RUN curl -sSL -O https://packages.microsoft.com/config/debian/12/packages-microsoft-prod.deb \
#  && dpkg -i packages-microsoft-prod.deb \
#  && rm packages-microsoft-prod.deb \
#  && apt-get update \
#  && apt-get install -y msodbcsql18 mssql-tools18 unixodbc-dev \
#  && rm -rf /var/lib/apt/lists/*
#  ENV PATH=$PATH:/opt/mssql-tools18/bin

# Oracle support.

# COPY --from=oracle /opt/oracle /opt/oracle
# ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/oracle/instantclient_23_26
# RUN ln -s /usr/lib/x86_64-linux-gnu/libaio.so.1t64 /usr/lib/x86_64-linux-gnu/libaio.so.1

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/init-database /app/init-database
COPY --from=senzingsdk /opt/senzing/er/resources/schema /opt/senzing/er/resources/schema

# Run as non-root container.

USER 1001

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
