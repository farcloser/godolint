# This file is the project's own — add recipes below. Keep the import: it
# mounts every shared limen task under `just do ...`.
import '.limen/just/main.just'

lint: do::lint::default do::lint::go::default do::lint::go::deadcode
fix: do::fix::default do::fix::go::default
test: do::test::go::unit do::test::go::race do::test::go::bench do::test::go::cover do::test::go::profile

# Regenerate from third-party/hadolint: DLxxxx metadata files are ALWAYS
# rewritten; stubs and ported tests only when absent. See internal/rules/generate.go.
generate:
    go generate ./internal/rules
