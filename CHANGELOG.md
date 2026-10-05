# Changelog

All notable changes are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow
[Semantic Versioning](https://semver.org/): the versioned surface is the
`godolint` command line and the `sdk` package.

## [Unreleased]

## [0.2.0] - 2026-10-05

### Changed

- The module path is `github.com/forkcloser/godolint`, after the repository
  moved from farcloser: `go install` and imports take the new path, since the
  Go toolchain refuses a module fetched under a path its go.mod does not
  declare.
- Logs go through `log/slog` on standard error, colored on a terminal and JSON
  when redirected. `--log-level` (debug, info, warn, error) sets the level
  and reads `LOG_LEVEL`; an unknown level is now a usage error, and the
  trace, fatal and panic spellings are gone.
- A run with no Dockerfile prints the usage with the error. The flags, the
  findings written as JSON on standard output and the exit status (1 on a
  violation or a failure) are unchanged.

### Added

- `--shellcheck-rcfile=FILE` hands a shellcheckrc to shellcheck when RUN
  instructions are checked (shellcheck 0.10.0 or later). In the `sdk`,
  `WithShellcheck` takes `ShellcheckOption`s, the first being
  `WithShellcheckRCFile(path)`; existing calls compile unchanged.

### Fixed

- Shellcheck's report is read from its standard output alone: a warning on
  its standard error no longer lands inside the JSON and fails the check.
