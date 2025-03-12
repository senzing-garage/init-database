# Makefile extensions for linux.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@nowhere/tmp/sqlite/G2C.db
# SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@/MYPRIVATE_DB?mode=memory&cache=shared
PATH := $(MAKEFILE_DIRECTORY)/bin:/$(HOME)/go/bin:$(PATH)

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: linux/amd64


.PHONY: clean-osarch-specific
clean-osarch-specific:
	@docker rm  --force $(DOCKER_CONTAINER_NAME) 2> /dev/null || true
	@docker rmi --force $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_IMAGE_NAME) $(DOCKER_SUT_IMAGE_NAME) 2> /dev/null || true
	@rm -f  $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -f  $(MAKEFILE_DIRECTORY)/.coverage || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.html || true
	@rm -f  $(MAKEFILE_DIRECTORY)/coverage.out || true
	@rm -f  $(MAKEFILE_DIRECTORY)/cover.out || true
	@rm -fr $(TARGET_DIRECTORY) || true
	@rm -fr /tmp/sqlite || true
	@pkill godoc || true
	@docker-compose -f docker-compose.test.yaml down 2> /dev/null || true


.PHONY: coverage-osarch-specific
coverage-osarch-specific: export SENZING_LOG_LEVEL=TRACE
coverage-osarch-specific:
	@go test -v -coverprofile=coverage.out -p 1 ./...
	@go tool cover -html="coverage.out" -o coverage.html
	@xdg-open $(MAKEFILE_DIRECTORY)/coverage.html


.PHONY: dependencies-for-development-osarch-specific
dependencies-for-development-osarch-specific:
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin latest


.PHONY: docker-build-osarch-specific
docker-build-osarch-specific:
	@$(DOCKER_BUILDKIT) docker build \
		--tag $(DOCKER_IMAGE_NAME) \
		--tag $(DOCKER_IMAGE_NAME):$(BUILD_VERSION) \
		.


.PHONY: documentation-osarch-specific
documentation-osarch-specific:
	@pkill godoc || true
	@godoc &
	@xdg-open http://localhost:6060


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	$(info Hello World, from linux.)


.PHONY: package-osarch-specific
package-osarch-specific: docker-build-package
	@mkdir -p $(TARGET_DIRECTORY) || true
	@CONTAINER_ID=$$(docker create $(DOCKER_BUILD_IMAGE_NAME)); \
	docker cp $$CONTAINER_ID:/output/. $(TARGET_DIRECTORY)/; \
	docker rm -v $$CONTAINER_ID


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@mkdir /tmp/sqlite
	@touch /tmp/sqlite/G2C.db
	docker-compose -f docker-compose.test.yaml up --detach


.PHONY: test-osarch-specific
test-osarch-specific:
	@echo "SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -tags "libsqlite3 linux" -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-mssql-osarch-specific
test-mssql-osarch-specific: export SENZING_TOOLS_DATABASE_URL=mssql://sa:Passw0rd@localhost:1433/G2/?TrustServerCertificate=yes
test-mssql-osarch-specific:
	@echo "SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-mysql-osarch-specific
test-mysql-osarch-specific: export SENZING_TOOLS_DATABASE_URL=mysql://mysql:mysql@127.0.0.1:3306/G2
test-mysql-osarch-specific:
	@echo "SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-oracle-osarch-specific
test-oracle-osarch-specific: export SENZING_TOOLS_DATABASE_URL=oci://pdbadmin:Passw0rd@localhost:1521/FREEPDB1
test-oracle-osarch-specific: export SENZING_TOOLS_SQL_FILE=$(MAKEFILE_DIRECTORY)/rootfs/opt/senzing/er/resources/schema/szcore-schema-oracle-create.sql
test-oracle-osarch-specific:
	@echo "SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-oracle-sys-osarch-specific
test-oracle-sys-osarch-specific: export SENZING_TOOLS_DATABASE_URL=oci://sys:Passw0rd@localhost:1521/FREE/?sysdba=true&noTimezoneCheck=true
test-oracle-sys-osarch-specific: export SENZING_TOOLS_SQL_FILE=$(MAKEFILE_DIRECTORY)/rootfs/opt/senzing/er/resources/schema/szcore-schema-oracle-create.sql
test-oracle-sys-osarch-specific:
	@echo "test-oracle-sys-osarch-specific:  SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt


.PHONY: test-postgresql-osarch-specific
test-postgresql-osarch-specific: export SENZING_TOOLS_DATABASE_URL=postgresql://postgres:postgres@localhost:5432/G2/?sslmode=disable
test-postgresql-osarch-specific:
	@echo "SENZING_TOOLS_DATABASE_URL: ${SENZING_TOOLS_DATABASE_URL}"
	@go test -json -v -p 1 ./... 2>&1 | tee /tmp/gotest.log | gotestfmt

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux
only-linux:
	$(info Only linux has this Makefile target.)
