# Project-specific just recipes go here.
# For example, you might want to define your own lint-all / test-all tasks
# (the aggregates CI runs; a dependency names a recipe, never a bare module):

lint-all: lint::default lint::go::default lint::go::bce lint::go::escape lint::go::deadcode
test-all: test::go::unit test::go::race test::go::bench test::go::cover test::go::profile

# Regenerate rule stubs and ported tests from the vendored hadolint sources
# (third-party/hadolint): gen-rules rewrites the DLxxxx metadata files (always)
# and stubs unimplemented rules; gen-tests ports the test corpus for
# implemented rules. See internal/rules/generate.go.
generate:
    go generate ./internal/rules
