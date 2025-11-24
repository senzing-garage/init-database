# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], [markdownlint],
and this project adheres to [Semantic Versioning].

## [Unreleased]

-

## [0.8.1] - 2025-11-25

### Added in 0.8.1

- Deduplicate database URLs

## [0.8.0] - 2025-11-18

### Added in 0.8.0

- Option to load Senzing TruthSet

## [0.7.22] - 2025-11-14

### Changed in 0.7.22

- Update dependencies

## [0.7.21] - 2025-09-04

### Changed in 0.7.21

- Admin release

## [0.7.20] - 2025-09-04

### Changed in 0.7.20

- Option `install-senzing-configuration` renamed to `install-senzing-er-configuration`

## [0.7.19] - 2025-08-28

### Changed in 0.7.19

- Move from beta to production

## [0.7.18] - 2025-07-18

### Changed in 0.7.18

- Update dependencies

## [0.7.17] - 2025-07-09

### Changed in 0.7.17

- Update dependencies

## [0.7.16] - 2025-07-02

### Changed in 0.7.16

- Update dependencies

## [0.7.15] - 2025-06-18

### Changed in 0.7.15

- Update dependencies

## [0.7.14] - 2025-06-12

### Changed in 0.7.14

- Update artifacts

## [0.7.13] - 2025-06-04

### Changed in 0.7.13

- Improve error handling

## [0.7.12] - 2025-05-15

### Changed in 0.7.12

- Update dependencies

## [0.7.11] - 2025-04-18

### Changed in 0.7.11

- Update dependencies

## [0.7.10] - 2025-04-15

### Changed in 0.7.10

- Update dependencies

## [0.7.9] - 2025-04-09

### Changed in 0.7.9

- Update dependencies

## [0.7.8] - 2025-03-28

### Added in 0.7.8

## [0.7.7] - 2025-03-12

### Added in 0.7.7

- Support for MS SQL
- Support for Oracle

## [0.7.6] - 2025-03-03

### Changed in 0.7.6

- Update dependencies

## [0.7.5] - 2025-02-28

### Changed in 0.7.5

- Update dependencies

## [0.7.4] - 2024-12-10

### Changed in 0.7.4

- Update dependencies

## [0.7.3] - 2024-11-14

### Changed in 0.7.3

- Support SQLite in-memory databases

## [0.7.2] - 2024-10-31

### Changed in 0.7.2

- Update dependencies

## [0.7.1] - 2024-09-13

### Changed in 0.7.1

- Update dependencies

## [0.7.0] - 2024-08-29

### Changed in 0.7.0

- Change from `g2` to `sz`/`er`

## [0.6.2] - 2024-08-05

### Changed in 0.6.2

- Update to `template-go`
- Updated dependencies

## [0.6.1] - 2024-07-08

### Changed in 0.6.1

- Updated dependencies

## [0.6.0] - 2024-05-10

### Changed in 0.6.0

- Migrate from `g2` to `sz`
- Update dependencies

## [0.5.2] - 2024-03-19

### Changed in 0.5.2

- Update `Dockerfile` to senzing/senzingapi-runtime:3.9.0
- Update dependencies

## [0.5.1] - 2024-01-30

### Changed in 0.5.1

- Update dependencies

## [0.5.0] - 2024-01-04

### Changed in 0.5.0

- Renamed module to `github.com/senzing-garage/init-database`
- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies

## [0.4.3] - 2023-12-06

### Changed in 0.4.3

- Update `Dockerfile` to senzing/senzingapi-runtime:3.8.0
- Update SQL files in `/opt/senzing/g2/resources/schema`
- Update `/opt/senzing/g2/resources/templates/g2config.json`
- GitHub Action to push to DockerHub and ECR
- Update dependencies

## [0.4.2] - 2023-11-02

### Changed in 0.4.2

- Update dependencies

## [0.4.1] - 2023-10-25

### Changed in 0.4.1

- Refactor to [template-go](https://github.com/senzing-garage/template-go)
- Update dependencies

## [0.4.0] - 2023-10-04

### Changed in 0.4.0

- Supports SenzingAPI 3.8.0
- Deprecated functions have been removed
- Update dependencies

## [0.3.2] - 2023-09-01

### Changed in 0.3.2

- Last version before SenzingAPI 3.8.0

## [0.3.1] - 2023-08-09

### Changed in 0.3.1

- Refactor to `template-go`
- Update dependencies

## [0.3.0] - 2023-07-21

### Added in 0.3.0

- Support for
  [SENZING_TOOLS_ENGINE_CONFIGURATION_FILE](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_engine_configuration_file) and
  [SENZING_TOOLS_SQL_FILE](https://github.com/senzing-garage/knowledge-base/blob/main/lists/environment-variables.md#senzing_tools_sql_file)

### Changed in 0.3.0

- Update `Dockerfile` to senzing/senzingapi-runtime:3.6.0
- Begin work on multi-platform support
- Update dependencies

## [0.2.6] - 2023-06-16

### Changed in 0.2.6

- Update dependencies

## [0.2.5] - 2023-05-26

### Changed in 0.2.5

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4

## [0.2.4] - 2023-05-17

### Changed in 0.2.4

- Added `gosec`
- Updated dependencies
  - github.com/senzing-garage/go-observing v0.2.5

## [0.2.3] - 2023-05-12

### Changed in 0.2.3

- In `Dockerfile`
  - `golang:1.20.4@sha256:31a8f92b17829b3ccddf0add184f18203acfd79ccc1bcb5c43803ab1c4836cca`
  - `senzing/senzingapi-runtime:3.5.2`
- Update dependencies

## [0.2.2] - 2023-04-21

### Changed in 0.2.2

- Update dependencies

## [0.2.1] - 2023-04-18

### Changed in 0.2.1

- Updated dependencies
- Migrated from `github.com/senzing-garage/go-logging/logger` to `github.com/senzing-garage/go-logging/logging`

## [0.2.0] - 2023-03-28

### Changed in 0.2.0

- Repository name change from `initdatabase` to `init-database`

## [0.1.8] - 2023-03-27

### Changed in 0.1.8

- Fixed configuration nil issue #28
- Last release before name change to `init-database`

## [0.1.7] - 2023-03-23

### Changed in 0.1.7

- Update dependencies
- If SQLite file doesn't exist, create it

## [0.1.6] - 2023-03-14

### Changed in 0.1.6

- Update dependencies
- Standardize use of Viper/Cobra

## [0.1.5] - 2023-03-13

### Fixed in 0.1.5

- Fixed issue silent error when connecting to database.

## [0.1.4] - 2023-03-08

### Fixed in 0.1.4

- Fixed issue with Cobra, Viper, and subcommand parameters, again.

## [0.1.3] - 2023-03-07

### Fixed in 0.1.3

- Fixed issue with Cobra, Viper, and subcommand parameters

## [0.1.2] - 2023-03-06

### Added to 0.1.2

- Fixed issue with `SENZING_TOOLS_DATABASE_URL`

## [0.1.1] - 2023-03-03

### Added to 0.1.1

- Organize input parameters

## [0.1.0] - 2023-03-03

### Added to 0.1.0

- Initial artifacts

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[markdownlint]: https://dlaa.me/markdownlint/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
