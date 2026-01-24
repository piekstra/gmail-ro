# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Binary renamed to `gmro`** - The CLI binary is now `gmro` (short for gmail-readonly). Install via `brew install gmail-readonly`, run with `gmro`. ([#66](https://github.com/open-cli-collective/gmail-ro/pull/66))
- Module path migrated to `github.com/open-cli-collective/gmail-ro` ([#63](https://github.com/open-cli-collective/gmail-ro/pull/63))

### Added

- `gmro init` command for guided OAuth setup with clear instructions ([#61](https://github.com/open-cli-collective/gmail-ro/pull/61))
- `gmro config show`, `config test`, `config clear` subcommands for credential management ([#61](https://github.com/open-cli-collective/gmail-ro/pull/61))
- Secure OAuth token storage in system keychain (macOS) or secret-tool (Linux) ([#59](https://github.com/open-cli-collective/gmail-ro/pull/59))

### Fixed

- Add nil check for message payload before accessing headers ([#55](https://github.com/open-cli-collective/gmail-ro/pull/55))
