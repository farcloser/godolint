# DO NOT EDIT MANUALLY.
# This file provides shared tasks common to all projects and managed by limen.
# If you want to add project specific tasks, you can do so in the `project.just` file.

set working-directory := '..'

# Silence just's per-line command echo; recipes announce themselves via _banner.
set quiet

# Shared private _banner recipe (see _lib.just).
import '_lib.just'

default: limen just yaml

# limen itself: repair what `just lint limen` reports — rewrite drifted
# canonical files, create missing mandatory ones, restore aqua pins and
# regenerate aqua-checksums.json (network, when pins changed). Same LIMEN_BIN
# override as the lint recipe (see lint.just).
limen: (_banner "fix" "limen")
    ${LIMEN_BIN:-limen} fix .

just: (_banner "fix" "just")
    #!/usr/bin/env bash
    set -euo pipefail
    # --cached --others --exclude-standard: tracked *and* new untracked files,
    # still honoring .gitignore — so a not-yet-staged file is formatted too.
    git ls-files -z --cached --others --exclude-standard 'justfile' 'Justfile' '*.just' '**/justfile' \
      | xargs -0 -I{} just --fmt --justfile {}

rust: (_banner "fix" "rust")
    cargo clippy --fix --all-targets --all-features --allow-dirty --allow-staged
    cargo fmt --all

# Go fixers live in their own submodule: `just fix go` runs them all,
# `just fix go <recipe>` (e.g. `just fix go mod`) runs one.
mod go 'fix-go.just'

yaml: (_banner "fix" "yaml")
    yamlfmt -conf .just/.yamlfmt .
