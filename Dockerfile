# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.22.3-bullseye@sha256:78b171fe51f25b8c3197710f281dacae94759a254e9a486576005f9dadba9e7d
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.10.1

# -----------------------------------------------------------------------------
# Stage: senzingapi_runtime
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as senzingapi_runtime

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2024-05-22
LABEL Name="senzing/init-database-builder" \
      Maintainer="support@senzing.com" \
      Version="0.5.2"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/init-database

# Copy files from prior stage.

COPY --from=senzingapi_runtime  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzingapi_runtime  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Build go program.

WORKDIR ${GOPATH}/src/init-database
RUN make build

# Copy binaries to /output.

RUN mkdir -p /output \
      && cp -R ${GOPATH}/src/init-database/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2024-05-22
LABEL Name="senzing/init-database" \
      Maintainer="support@senzing.com" \
      Version="0.5.2"

# Copy local files from the Git repository.

COPY ./rootfs /

# Copy files from prior stage.

COPY --from=go_builder "/output/linux/init-database" "/app/init-database"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
