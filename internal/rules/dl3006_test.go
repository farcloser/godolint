package rules

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
)

// Auto-generated tests for DL3006 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3006Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3006(t *testing.T) {
	allRules := []rule.Rule{DL3006()}

	t.Run("local aliases are OK to be untagged", func(t *testing.T) {
		dockerfile := `FROM golang:1.9.3-alpine3.7 AS build
RUN foo
FROM build as unit-test
RUN bar
FROM alpine:3.7
RUN baz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")
	})

	t.Run("no untagged", func(t *testing.T) {
		dockerfile := `FROM debian`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3006")
	})

	t.Run("no untagged with name", func(t *testing.T) {
		dockerfile := `FROM debian AS builder`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3006")
	})

	t.Run("other untagged cases are not ok", func(t *testing.T) {
		dockerfile := `FROM golang:1.9.3-alpine3.7 AS build
RUN foo
FROM node as unit-test
RUN bar
FROM alpine:3.7
RUN baz`
		violations := LintDockerfile(dockerfile, allRules)

		AssertContainsViolation(t, violations, "DL3006")
	})

	t.Run("scratch", func(t *testing.T) {
		dockerfile := `FROM scratch`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")
	})

	t.Run("untagged digest is not an error", func(t *testing.T) {
		dockerfile := `FROM ruby@sha256:f1dbca0f5dbc9`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")
	})

	t.Run("untagged digest is not an error", func(t *testing.T) {
		dockerfile := `FROM ruby:2`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")
	})

	t.Run("using args is not an error", func(t *testing.T) {
		dockerfile := `FROM ${VALUE}`
		violations := LintDockerfile(dockerfile, allRules)

		AssertNoViolation(t, violations, "DL3006")
	})
}
