# Changelog

All notable changes to this project are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project uses [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.4.0] - 2026-09-13

### Added
- Container image, published to `ghcr.io/gllm-dev/molinets` for linux/amd64 and linux/arm64 on every GitHub release.
- GitHub Actions workflow: `go mod tidy` and Wire output checks, vet, and race-enabled tests.
- README with quick start, client setup, tool and notation reference.

## [0.3.0] - 2026-08-16

### Added
- `stdio` command and transport, for MCP clients that run the server as a local process. Logs go to stderr so stdout carries only the protocol.

## [0.2.0] - 2026-08-02

### Added
- Application layer:
  - Grinder catalog: list, resolve a free-text name, report ambiguous names, deduplicate targets.
  - Conversion use cases: setting to setting, setting range for a brew method, microns to setting.
- MCP server over Streamable HTTP, stateless with JSON responses, on `/mcp`, with `/healthz`.
- Tools: `list_grinders`, `convert_grind_setting`, `grind_setting_for_method` and `microns_to_setting`, all read-only.
- Stable error codes on failed tool calls (`GRINDER_NOT_FOUND`, `AMBIGUOUS_GRINDER`, `INVALID_SETTING`, `SETTING_OUT_OF_RANGE`, `METHOD_NOT_FOUND`, `METHOD_UNREACHABLE`, `INVALID_ARGUMENT`).
- Warning when a source setting is at burr contact.
- `http` command (cobra), `PORT` configuration, and Wire dependency injection.
- Tests: application tests against hand-written mocks of the ports, and end-to-end tests with the SDK client.

## [0.1.1] - 2026-07-18

### Added
- Repository ports and embedded repositories for grinders and brew methods. Grinder names resolve by id, name or alias (accents and punctuation ignored). Brew methods also accept brewer names such as Chemex, Origami, Switch, Clever or Moccamaster.
- Dataset tests:
  - 263 reference points from the Honest Coffee Guide tables,
  - clean formatting and parse round trips across every grinder's range,
  - monotonic calibrations,
  - exact round trips between every pair of grinders.

## [0.1.0] - 2026-07-05

### Added
- Domain model:
  - grinders with a native notation (clicks, rotation.number.click, rotation.clicks, number.click, thirds, decimal dial),
  - a piecewise linear calibration to microns that never goes below zero,
  - a usable particle-size range, never finer than 40 µm.
- Strict notation parsing that rejects impossible settings with the valid range.
- Brew methods with their particle-size ranges.
- Dataset of 34 manual and electric grinders.

[0.4.0]: https://github.com/gllm-dev/molinets/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/gllm-dev/molinets/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/gllm-dev/molinets/compare/v0.1.1...v0.2.0
[0.1.1]: https://github.com/gllm-dev/molinets/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/gllm-dev/molinets/releases/tag/v0.1.0
