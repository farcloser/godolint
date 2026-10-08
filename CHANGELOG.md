# Changelog

All notable changes are recorded here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and versions follow
[Semantic Versioning](https://semver.org/): the versioned surface is the
`godolint` command line and the `sdk` package.

## [Unreleased]

## [0.3.0] - 2026-10-08

### Added

- Six rules, bringing the port to hadolint 2.15.1's 71: DL3063 (a stage
  name is not a reserved word), DL3064 (ARG and ENV carry no sensitive
  data), DL3065 (`FROM --platform` is not the default TARGETPLATFORM),
  DL3066 (a USER is a numeric id or a declared ARG), DL3067 (a stage's
  whole filesystem is not copied), and DL3056 (a label the schema types as
  a version is a semantic version), which 2.14.0 had and the generator
  could not read.

### Changed

- The rules are generated and tested from hadolint v2.15.1, pinned in
  `pins.yaml` with its digest; CI regenerates from the pin and diffs.
- Ten rules follow upstream's current behaviour: DL3033 and DL3041 (a
  version in a variable is pinned when the variable is defined, ENV per
  stage through a FROM alias, ARG across the file; DL3041 skips `group`),
  DL3013 (pipx counts as pip; `--python` and `--root-user-action` take a
  value), DL3016 (`--registry` takes a value), DL3018 (a fuzzy apk
  constraint is a pin), DL3062 (a local path has no version; `go run`'s
  arguments are the program's), DL3032, DL3036, DL3040 and DL3060 (a clean
  counts only after the install), DL3040 (dnf's other install spellings,
  for dnf and microdnf).
- A plain variable reference in a shell word keeps its name (`$v` and
  `${v}` read as `${v}`; operators and other expansions stay masked), as
  hadolint's simplify does, so rules can tell which variable a word uses.
- A failure of a rule whose severity is Ignore (DL3057) is reported by the
  rule and dropped by the linter and the command, where hadolint drops it.
- Built on mycophonic/primordium v0.12.0.

### Fixed

- Sixteen rule messages (DL3001, DL3004, DL3007, DL3008, DL3013, DL3016,
  DL3018, DL3019, DL3027, DL3028, DL3042, DL3047, DL3061, DL4003, DL4004,
  DL4006) carried a run of 7 to 15 spaces where hadolint's has one: the
  generator kept the indentation inside hadolint's Haskell string gaps,
  which the Haskell compiler deletes. The messages now read as hadolint
  emits them, and the generator deletes the gap. DL1001's two spaces are
  hadolint's own (one on each side of its gap) and stay. Reported by The
  Designer from the site's hero sample.

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
