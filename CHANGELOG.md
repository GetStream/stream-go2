# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [8.7.0](https://github.com/GetStream/stream-go2/compare/v8.6.0...v8.7.0) (2025-03-19)

## [8.6.0](https://github.com/GetStream/stream-go2/compare/v8.5.0...v8.6.0) (2025-03-07)

## [8.5.0](https://github.com/GetStream/stream-go2/compare/v8.4.3...v8.5.0) (2025-03-05)

### [8.4.3](https://github.com/GetStream/stream-go2/compare/v8.4.2...v8.4.3) (2024-06-25)

### [8.4.2](https://github.com/GetStream/stream-go2/compare/v8.4.1...v8.4.2) (2024-06-18)

### [8.4.1](https://github.com/GetStream/stream-go2/compare/v8.4.0...v8.4.1) (2024-06-14)

## [8.4.0](https://github.com/GetStream/stream-go2/compare/v8.3.0...v8.4.0) (2024-01-09)

## [8.3.0](https://github.com/GetStream/stream-go2/compare/v8.2.1...v8.3.0) (2023-11-20)

### [8.2.1](https://github.com/GetStream/stream-go2/compare/v8.2.0...v8.2.1) (2023-08-17)

## [8.2.0](https://github.com/GetStream/stream-go2/compare/v8.1.0...v8.2.0) (2023-08-17)

## [8.1.0](https://github.com/GetStream/stream-go2/compare/v8.0.2...v8.1.0) (2023-07-25)

### [8.0.2](https://github.com/GetStream/stream-go2/compare/v8.0.1...v8.0.2) (2023-05-10)

### [8.0.1](https://github.com/GetStream/stream-go2/compare/v8.0.0...v8.0.1) (2023-02-13)


### Bug Fixes

* link to correct repo ([c771c1f](https://github.com/GetStream/stream-go2/commit/c771c1fe49c1ae1ef502fd3015383effe0bbc317))
* use v8 in tests ([a9406ad](https://github.com/GetStream/stream-go2/commit/a9406adb46678089d6e957299efcfea134494334))

## [8.0.0](https://github.com/GetStream/stream-go2/compare/v7.1.0...v8.0.0) (2022-10-20)


### Features

* add user enrichment for reactions ([#137](https://github.com/GetStream/stream-go2/issues/137)) ([d88c659](https://github.com/GetStream/stream-go2/commit/d88c659dd5520cdd9bc8388912857834f0b4086b))

## [7.1.0](https://github.com/GetStream/stream-go2/compare/v7.0.1...v7.1.0) (2022-10-04)


### Features

* add time fields to reactions ([#134](https://github.com/GetStream/stream-go2/issues/134)) ([bd5966c](https://github.com/GetStream/stream-go2/commit/bd5966c3eb5930cd050844412fe093060ad64222))

### [7.0.1](https://github.com/GetStream/stream-go2/compare/v7.0.0...v7.0.1) (2022-06-21)


### ⚠ BREAKING CHANGES

* rename session token methods (#126)

### Features

* **go_version:** bump to v1.17 ([#125](https://github.com/GetStream/stream-go2/issues/125)) ([0c1e87c](https://github.com/GetStream/stream-go2/commit/0c1e87c0451859787d95de11a955253d8ee00b49))


* rename session token methods ([#126](https://github.com/GetStream/stream-go2/issues/126)) ([39fbcf7](https://github.com/GetStream/stream-go2/commit/39fbcf75c16aa26c70c12afbd5d4d9faab8d5a4e))

## [7.0.0](https://github.com/GetStream/stream-go2/compare/v6.4.2...v7.0.0) (2022-05-10)


### Features

* **context:** add context as first argument ([#123](https://github.com/GetStream/stream-go2/issues/123)) ([9612a24](https://github.com/GetStream/stream-go2/commit/9612a24b921d4aeb8ab4b22e8c5ddd93e84ecf9e))

## [6.4.2] 2022-03-10

- Improve keep-alive settings of the default client.

## [6.4.1] 2022-03-09

- Handle activity references in foreign id for enrichment. Enriched activity is put into `foreign_id_ref` under `Extra`.

## [6.4.0] 2021-12-15

- Add new flags for reaction pagination
- Fix parsing next url in reaction pagination

## [6.3.0] 2021-12-03

- Add new reaction flags
  - first reactions
  - reaction count
  - own children kind filter

## [6.2.0] 2021-11-19

- Add user id support into reaction own children for filtering

## [6.1.0] 2021-11-15

- Expose created_at/updated_at in groups for aggregated/notification feeds

## [6.0.0] 2021-11-12

- Add enrichment options into read activity endpoints
- Move support into go 1.16 & 1.17

## [5.7.2] 2021-08-04

- Dependency upgrade for unmaintained jwt

## [5.7.1] 2021-07-01

- Fix godoc issues

## [5.7.0] 2021-06-04

- Add follow stats endpoint support ([#108](https://github.com/GetStream/stream-go2/pull/108))
- Run CI with 1.15 and 1.16
- Add a changelog to document changes
