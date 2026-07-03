# DO NOT EDIT MANUALLY.
# This file provides shared tasks common to all projects and managed by limen.
# If you want to add project specific tasks, you can do so in the `project.just` file.

# Imported flat into the canonical Justfile — not a module — so the recipe can
# take arguments (`just release v1.2.3`; a module invocation would parse the
# tag as a recipe path). Deliberately no settings and no _lib import here: an
# import's settings would apply to the whole root justfile.

# Release, in two lanes that share every guard (see book/recipes.md):
#
#   just release --cut vX.Y.Z   the CI lane's local half: verify a clean tree,
#                               create the signed tag (your key signs the
#                               INTENT), push — the release workflow builds,
#                               signs (keyless cosign), and publishes.
#   just release --ci           the CI lane's workflow half (GitHub Actions
#                               only): the pushed tag triggered us; run
#                               goreleaser with keyless signing.
#   just release vX.Y.Z         the fully local lane: tag + goreleaser +
#                               key-based cosign from this machine. Requires
#                               GITHUB_TOKEN, COSIGN_KEY, COSIGN_PASSWORD.
#                               The lane for private repos (nothing touches
#                               Rekor's public log) — and the escape hatch
#                               when CI is down.
#   just release --dry-run      unsigned local snapshot into build/release/:
#                               no tag, no token, no publish — and works on a
#                               dirty tree, by design.
#
# Tag pushes are always exactly one tag, never --tags: stray local tags stay
# local. A tag already on HEAD is reused, so a failed publish retries safely.
[doc('Release: `--cut vX.Y.Z` tag for CI · `vX.Y.Z` fully local · `--ci` (workflow) · `--dry-run`')]
release *args:
    #!/usr/bin/env bash
    set -euo pipefail
    # Releasing is opt-in, by carrying a goreleaser config (project-owned, like
    # project.just — limen neither pins nor seeds it). Checked before anything
    # else: past this point the recipe creates tags and pushes.
    if [ ! -f .goreleaser.yaml ] && [ ! -f .goreleaser.yml ]; then
        echo "this project has no .goreleaser.yaml — goreleaser releases are opt-in (see book/recipes.md)." >&2
        exit 1
    fi
    usage() {
        echo "usage: just release --cut <tag>   # signed tag + push; CI builds, signs (keyless), publishes" >&2
        echo "       just release <tag>         # fully local release (key-based cosign)" >&2
        echo "       just release --dry-run     # local unsigned snapshot into build/release/" >&2
        echo "       just release --ci          # workflow half of --cut (GitHub Actions only)" >&2
    }
    require_clean() {
        # A dirty tree is the most common failure and goreleaser hard-refuses
        # it — check long before anything is created or pushed.
        if [ -n "$(git status --porcelain)" ]; then
            echo "the working tree is dirty — commit or stash before releasing:" >&2
            git status --short >&2
            exit 1
        fi
    }
    require_version_tag() {
        if ! printf '%s' "$1" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+(-[0-9A-Za-z][0-9A-Za-z.-]*)?$'; then
            echo "'$1' is not a version tag (vX.Y.Z or vX.Y.Z-prerelease, e.g. v0.0.0-test.1)." >&2
            exit 2
        fi
    }
    tag_and_push() {
        # Create the signed tag — unless it already exists pointing at HEAD, so
        # a retry after a failed publish resumes instead of failing here.
        if existing=$(git rev-parse -q --verify "refs/tags/$1^{commit}"); then
            if [ "$existing" != "$(git rev-parse HEAD)" ]; then
                echo "tag $1 already exists and does not point at HEAD — refusing." >&2
                exit 1
            fi
            echo "tag $1 already on HEAD — reusing it."
        else
            git tag -s "$1" -m "$1"
        fi
        git push origin HEAD
        git push origin "refs/tags/$1"
    }
    set -- {{ args }}
    case "${1:-}" in
        --dry-run)
            # Snapshot builds skip signing: no key material needed for a dry-run.
            exec goreleaser release --snapshot --clean --skip=sign
            ;;
        --cut)
            [ -n "${2:-}" ] || { usage; exit 2; }
            require_clean
            require_version_tag "$2"
            tag_and_push "$2"
            echo "tag $2 pushed — the release workflow builds, signs, and publishes from here."
            ;;
        --ci)
            # The workflow half: the tag exists (its push triggered this run,
            # and CI checks out exactly that commit), publishing uses the
            # workflow's GITHUB_TOKEN, and signing is keyless — the OIDC
            # identity replaces any key.
            if [ -z "${GITHUB_ACTIONS:-}" ]; then
                echo "--ci is the workflow half — locally, use 'just release --cut <tag>' (or a fully local 'just release <tag>')." >&2
                exit 1
            fi
            if [ -z "${ACTIONS_ID_TOKEN_REQUEST_URL:-}" ]; then
                echo "no OIDC token available — the workflow needs 'id-token: write' for keyless signing." >&2
                exit 1
            fi
            if [ -z "${GITHUB_TOKEN:-}" ]; then
                echo "GITHUB_TOKEN is not set — pass the workflow token to this step's env." >&2
                exit 1
            fi
            if ! git describe --exact-match --tags HEAD >/dev/null 2>&1; then
                echo "HEAD carries no tag — this lane only runs on a tag push." >&2
                exit 1
            fi
            # Keyless is enforced, not defaulted: a key leaking into the CI
            # environment must not silently flip the signing mode.
            unset COSIGN_KEY
            exec goreleaser release --clean
            ;;
        "")
            usage
            exit 2
            ;;
        *)
            tag="$1"
            require_clean
            require_version_tag "$tag"
            if [ -z "${GITHUB_TOKEN:-}" ]; then
                echo "GITHUB_TOKEN is not set — goreleaser needs it to publish the GitHub release." >&2
                exit 1
            fi
            if [ -z "${COSIGN_KEY:-}" ] || [ ! -f "${COSIGN_KEY}" ]; then
                echo "COSIGN_KEY is not set or does not point at the cosign private key —" >&2
                echo "local releases are key-signed (see .goreleaser.yaml). Generate a key pair once" >&2
                echo "with 'cosign generate-key-pair', commit cosign.pub, keep cosign.key private," >&2
                echo "then export COSIGN_KEY=/path/to/cosign.key (COSIGN_PASSWORD is prompted or exported)." >&2
                echo "(Or use the CI lane instead: 'just release --cut $tag' — no key needed.)" >&2
                exit 1
            fi
            tag_and_push "$tag"
            exec goreleaser release --clean
            ;;
    esac
