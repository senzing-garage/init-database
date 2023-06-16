# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
[markdownlint](https://dlaa.me/markdownlint/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

-

## [0.2.6] - 2023-06-16

### Changed in 0.2.6

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.6
  - github.com/senzing/go-common v0.1.4
  - github.com/senzing/go-databasing v0.2.5
  - github.com/senzing/go-logging v1.2.6
  - github.com/senzing/go-observing v0.2.6
  - github.com/senzing/go-sdk-abstract-factory v0.3.0
  - github.com/senzing/senzing-tools v0.2.9-0.20230613173043-18f1bd4cafdb
  - github.com/spf13/viper v1.16.0
  - google.golang.org/grpc v1.56.0

## [0.2.5] - 2023-05-26

### Changed in 0.2.5

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.4

## [0.2.4] - 2023-05-17

### Changed in 0.2.4

- Added `gosec`
- Updated dependencies
  - github.com/senzing/go-observing v0.2.5

## [0.2.3] - 2023-05-12

### Changed in 0.2.3

- In `Dockerfile`
  - `golang:1.20.4@sha256:31a8f92b17829b3ccddf0add184f18203acfd79ccc1bcb5c43803ab1c4836cca`
  - `senzing/senzingapi-runtime:3.5.2`
- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.2
  - github.com/senzing/go-common v0.1.3
  - github.com/senzing/go-databasing v0.2.4
  - github.com/senzing/go-logging v1.2.3
  - github.com/senzing/go-observing v0.2.3
  - github.com/senzing/go-sdk-abstract-factory v0.2.3
  - github.com/senzing/senzing-tools v0.2.4
  - google.golang.org/grpc v1.55.0

## [0.2.2] - 2023-04-21

### Changed in 0.2.2

- Update dependencies
  - github.com/senzing/g2-sdk-go v0.6.1
  - github.com/senzing/go-sdk-abstract-factory v0.2.2

## [0.2.1] - 2023-04-18

### Changed in 0.2.1

- Updated dependencies
- Migrated from `github.com/senzing/go-logging/logger` to `github.com/senzing/go-logging/logging`

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
- If Sqlite file doesn't exist, create it

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
