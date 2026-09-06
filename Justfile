# This file is the project's own — add recipes below. Keep the import: it
# mounts every shared limen task under `just do ...`.
import '.limen/just/main.just'

lint: do::lint::default do::lint::go::default do::lint::go::deadcode
fix: do::fix::default do::fix::go::default
test: do::test::go::unit do::test::go::race do::test::go::bench do::test::go::cover do::test::go::profile

# Regenerate rule stubs and ported tests from the vendored hadolint sources
# (third-party/hadolint): gen-rules rewrites the DLxxxx metadata files (always)
# and stubs unimplemented rules; gen-tests ports the test corpus for
# implemented rules. See internal/rules/generate.go.
generate:
    go generate ./internal/rules
