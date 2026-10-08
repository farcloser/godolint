# Changelog

All notable changes are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow
[Semantic Versioning](https://semver.org/): the versioned surface is the
`godolint` command line and the `sdk` package.

## [Unreleased]

## [0.2.1] - 2026-10-08

### Changed

- Built on mycophonic/primordium v0.12.0.

### Fixed

- Seventeen rule messages (DL1001, DL3001, DL3004, DL3007, DL3008, DL3013,
  DL3016, DL3018, DL3019, DL3027, DL3028, DL3042, DL3047, DL3061, DL4003,
  DL4004, DL4006) carried a run of 7 to 15 spaces where hadolint's has one:
  the generator kept the indentation inside hadolint's Haskell string gaps,
  which the Haskell compiler deletes. The messages now read as hadolint
  emits them, and the generator collapses the gap. Reported by The Designer
  from the site's hero sample.

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
