# DO NOT EDIT MANUALLY.
# This file provides shared tasks common to all projects and managed by limen.
# If you want to add project specific tasks, you can do so in the `project.just` file.

set working-directory := '..'

# Silence just's per-line command echo; recipes announce themselves via _banner.
set quiet

# Shared private _banner recipe (see _lib.just).
import '_lib.just'

rust: (_banner "lint" "rust")
    cargo fmt --all --check
    cargo clippy --all-targets --all-features -- --deny warnings
