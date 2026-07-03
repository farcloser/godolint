# DO NOT EDIT MANUALLY.
# This file provides shared tasks common to all projects and managed by limen.
# If you want to add project specific tasks, you can do so in the `project.just` file.

# Go tests live in their own submodule: `just test go` runs the unit tests,
# `just test go <recipe>` (e.g. `just test go race`) runs one task.
#
# Bare `just test` deliberately refuses: tests are always language-bound, so
# there is nothing a universal default could run that would not fail (or lie)
# on some repo. Which suites apply is each project's call — name one
# explicitly, and define the project-wide aggregate as `test-all` in
# project.just (the CI workflow's entry point, mirroring lint-all).
default:
    @echo 'just test has no default — name a suite (e.g. `just test go`), or run the project'"'"'s `just test-all`.' >&2
    @exit 1

mod go 'test-go.just'
