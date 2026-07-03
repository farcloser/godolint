# DO NOT EDIT MANUALLY.
# This file provides shared tasks common to all projects and managed by limen.
# If you want to add project specific tasks, you can do so in the `project.just` file.

set working-directory := '..'

# Silence just's per-line command echo; recipes announce themselves via _banner.
set quiet

# Go builds live in their own submodule: `just build go` makes the release
# binaries, `just build go <recipe>` (debug, race, static) makes a variant.
mod go 'build-go.just'

# A nested submodule's default is not addressable as a dependency path, so the
# bare `just build` delegates through a subprocess (same pattern as lint-all).
default:
    just build go
