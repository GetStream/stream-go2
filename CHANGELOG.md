# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
