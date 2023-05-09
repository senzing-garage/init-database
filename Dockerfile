# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.20.4@sha256:31a8f92b17829b3ccddf0add184f18203acfd79ccc1bcb5c43803ab1c4836cca
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.5.2

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-05-09
LABEL Name="senzing/init-database-builder" \
      Maintainer="support@senzing.com" \
      Version="0.1.4"

# Build arguments.

ARG PROGRAM_NAME="unknown"
ARG BUILD_VERSION=0.0.0
ARG BUILD_ITERATION=0
ARG GO_PACKAGE_NAME="unknown"

# Copy remote files from DockerHub.

COPY --from=senzing/senzingapi-runtime:3.4.2  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzing/senzingapi-runtime:3.4.2  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Copy local files from the Git repository.

COPY . ${GOPATH}/src/${GO_PACKAGE_NAME}

# Build go program.

WORKDIR ${GOPATH}/src/${GO_PACKAGE_NAME}
RUN make build

# --- Test go program ---------------------------------------------------------

# Run unit tests.

# RUN go get github.com/jstemmer/go-junit-report \
#  && mkdir -p /output/go-junit-report \
#  && go test -v ${GO_PACKAGE_NAME}/... | go-junit-report > /output/go-junit-report/test-report.xml

# Copy binaries to /output.

RUN mkdir -p /output \
      && cp -R ${GOPATH}/src/${GO_PACKAGE_NAME}/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2023-05-09
LABEL Name="senzing/init-database" \
      Maintainer="support@senzing.com" \
      Version="0.1.4"

# Copy files from repository.

COPY ./rootfs /

# Copy files from prior step.

COPY --from=go_builder "/output/linux/init-database" "/app/init-database"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/init-database"]
