# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.25.3-bookworm@sha256:51b6b12427dc03451c24f7fc996c43a20e8a8e56f0849dd0db6ff6e9225cc892
# ARG IMAGE_FINAL=senzing/senzingsdk-runtime:4.0.0@sha256:332d2ff9f00091a6d57b5b469cc60fd7dc9d0265e83d0e8c9e5296541d32a4aa
ARG IMAGE_FINAL=senzing/senzingsdk-runtime:latest

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
# Stage: mysql
# -----------------------------------------------------------------------------

# FROM ${IMAGE_FINAL} AS mysql
# ENV REFRESHED_AT=2024-08-01

# RUN mkdir /tmp/mysql \
#  && cd /tmp/mysql \
#  && wget https://cdn.mysql.com//Downloads/MySQL-9.5/mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar \
#  && tar -xf mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar

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
      # libssl3 \
      # libaio1 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

#  && ln -s /usr/lib/x86_64-linux-gnu/libmysqlclient.so  /opt/senzing/er/lib/libmysqlclient.so.21


# MySQL support

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
#  && dpkg -i --ignore-depends=libssl1.1 libmysqlclient21



# Install libssl1.1, libmysqlclient21, and mysql-connector-odbc
# RUN echo "deb http://deb.debian.org/debian-security/ bullseye-security main" |  tee /etc/apt/sources.list.d/bullseye-security.list
# RUN apt-get update
# RUN apt-get install libssl1.1
# RUN apt-get install libmysqlclient21
# RUN apt-get install mysql-connector-odbc


# RUN wget https://dev.mysql.com/get/Downloads/Connector-ODBC/8.0/mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#   && wget https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-common_8.0.20-1debian10_amd64.deb \
#   && wget http://repo.mysql.com/apt/debian/pool/mysql-8.0/m/mysql-community/libmysqlclient21_8.0.20-1debian10_amd64.deb \
#   && apt-get update \
#   && apt-get -y install \
#   ./mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#   ./mysql-common_8.0.20-1debian10_amd64.deb \
#   ./libmysqlclient21_8.0.20-1debian10_amd64.deb \
#   && rm \
#   ./mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#   ./mysql-common_8.0.20-1debian10_amd64.deb \
#   ./libmysqlclient21_8.0.20-1debian10_amd64.deb \
#   && rm -rf /var/lib/apt/lists/*



# COPY --from=mysql /tmp/mysql/libmysqlclient24_9.5.0-1debian13_amd64.deb /tmp/libmysqlclient24_9.5.0-1debian13_amd64.deb
# COPY --from=mysql /tmp/mysql/mysql-common_9.5.0-1debian13_amd64.deb     /tmp/mysql-common_9.5.0-1debian13_amd64.deb

# RUN wget https://dev.mysql.com/get/Downloads/Connector-ODBC/9.5/mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#  && apt-get update \
#  && apt-get -y install \
#       ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#       ./mysql-common_9.5.0-1debian13_amd64.deb \
#       ./libmysqlclient24_8.4.7-1debian13_amd64.deb \
#  && rm \
#       ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#       ./mysql-common_9.5.0-1debian13_amd64.deb \
#       ./libmysqlclient24_8.4.7-1debian13_amd64.deb \
#  && rm -rf /var/lib/apt/lists/*


#  RUN wget http://archive.ubuntu.com/ubuntu/pool/main/o/openssl/libssl1.1_1.1.0g-2ubuntu4_amd64.deb \
#   && dpkg -i ./libssl1.1_1.1.0g-2ubuntu4_amd64.deb \
#   && rm -rf ./libssl1.1_1.1.0g-2ubuntu4_amd64.deb


# RUN wget https://dev.mysql.com/get/mysql-apt-config_0.8.36-1_all.deb
# RUN dpkg-deb -c mysql-apt-config_0.8.36-1_all.deb
# RUN dpkg -i mysql-apt-config_0.8.36-1_all.deb
# RUN apt-get update
# RUN apt-get -y --no-install-recommends install mysql-client mysql-community-client
# RUN mysql --version


# RUN wget https://dev.mysql.com/get/Downloads/Connector-ODBC/9.5/mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#  && wget https://dev.mysql.com/get/Downloads/MySQL-9.5/mysql-common_9.5.0-1debian13_amd64.deb \
#  && wget https://repo.mysql.com/apt/debian/pool/mysql-9.5/m/mysql-community/libmysqlclient24_9.5.0-1debian13_amd64.deb \
#  && apt-get update \
#  && apt-get -y install \
#       ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#       ./mysql-common_9.5.0-1debian13_amd64.deb \
#       ./libmysqlclient24_8.4.7-1debian13_amd64.deb \
#  && rm \
#       ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
#       ./mysql-common_9.5.0-1debian13_amd64.deb \
#       ./libmysqlclient24_8.4.7-1debian13_amd64.deb \
#  && rm -rf /var/lib/apt/lists/*

# RUN wget https://dev.mysql.com/get/Downloads/Connector-ODBC/8.0/mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#  && wget https://dev.mysql.com/get/Downloads/MySQL-8.0/mysql-common_8.0.20-1debian10_amd64.deb \
#  && wget http://repo.mysql.com/apt/debian/pool/mysql-8.0/m/mysql-community/libmysqlclient21_8.0.20-1debian10_amd64.deb \
#  && apt-get update \
#  && apt-get -y install \
#       ./mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#       ./mysql-common_8.0.20-1debian10_amd64.deb \
#       ./libmysqlclient21_8.0.20-1debian10_amd64.deb \
#  && rm \
#       ./mysql-connector-odbc_8.0.20-1debian10_amd64.deb \
#       ./mysql-common_8.0.20-1debian10_amd64.deb \
#       ./libmysqlclient21_8.0.20-1debian10_amd64.deb \
# && rm -rf /var/lib/apt/lists/*

# RUN wget https://cdn.mysql.com/Downloads/MySQL-9.5/mysql-community-client_9.5.0-1debian13_amd64.deb


# RUN wget https://cdn.mysql.com//Downloads/MySQL-9.5/mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar \
#  && tar -xf mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar \
#  && apt-get -y install \
#       ./libmysqlclient24_9.5.0-1debian13_amd64.deb \
#       ./libmysqlclient-dev_9.5.0-1debian13_amd64.deb \
#       ./mysql-client_9.5.0-1debian13_amd64.deb \
#       ./mysql-common_9.5.0-1debian13_amd64.deb \
#       ./mysql-community-client_9.5.0-1debian13_amd64.deb \
#       ./mysql-community-client-core_9.5.0-1debian13_amd64.deb \
#       ./mysql-community-client-plugins_9.5.0-1debian13_amd64.deb

# RUN wget https://cdn.mysql.com//Downloads/MySQL-9.5/mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar \
#  && tar -xf mysql-server_9.5.0-1debian13_amd64.deb-bundle.tar \
#  && apt-get -y install \
#     ./libmysqlclient24_9.5.0-1debian13_amd64.deb

      # ./mysql-connector-odbc_9.5.0-1debian13_amd64.deb \
      # ./mysql-common_9.5.0-1debian13_amd64.deb \
      # ./libmysqlclient24_9.5.0-1debian13_amd64.deb



# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/init-database /app/init-database
COPY --from=senzingsdk /opt/senzing/er/resources/schema /opt/senzing/er/resources/schema
# COPY --from=oracle /opt/oracle /opt/oracle

# Run as non-root container

USER 1001

# Runtime environment variables.

# ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/opt/oracle/instantclient_23_9

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
