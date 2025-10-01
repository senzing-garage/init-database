# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_BUILDER=golang:1.25.1-bookworm@sha256:2960a1db140a9a6dd42b15831ec6f8da0c880df98930411194cf11875d433021
ARG IMAGE_FINAL=senzing/senzingsdk-runtime:4.0.0@sha256:332d2ff9f00091a6d57b5b469cc60fd7dc9d0265e83d0e8c9e5296541d32a4aa

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
 && apt-get -y install \
      libsqlite3-dev \
      wget

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

RUN make build

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

RUN apt-get update -qqq \
 && apt-get -yqqq install ${SENZING_APT_INSTALL_SETUP_PACKAGE}

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
 && apt-get -y install \
      libsqlite3-dev \
      # libaio1 \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=builder /output/linux/init-database /app/init-database
COPY --from=senzingsdk /opt/senzing/er/resources/schema /opt/senzing/er/resources/schema
# COPY --from=oracle /opt/oracle /opt/oracle

# Run as non-root container

USER 1001

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/er/lib/:/opt/oracle/instantclient_23_9:$LD_LIBRARY_PATH

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
