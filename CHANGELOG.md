# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial implementation of entropy check in Go
- Support for configurable warning and critical thresholds
- Command-line flags and environment variable configuration
- Metrics output in Graphite format via `--metrics` flag on check-entropy
- Standalone `metrics-entropy` command for dedicated metrics collection
- Configurable metric naming scheme
- Compatible with Sensu Go plugin SDK
- Multi-platform builds via GoReleaser

[Unreleased]: https://github.com/nmollerup/sensu-check-entropy/compare/v0.1.0...HEAD
