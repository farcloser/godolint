package rules_test

import (
	"testing"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
	"github.com/farcloser/godolint/internal/testutils"
)

// Auto-generated tests for DL3006 ported from hadolint test suite.
// Source: hadolint/test/Hadolint/Rule/DL3006Spec.hs
//
// To regenerate: go generate ./internal/rules

func TestDL3006(t *testing.T) {
	t.Parallel()

	allRules := []rule.Rule{
		rules.DL3006(),
	}

	t.Run(
		"local aliases are OK to be untagged",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM golang:1.9.3-alpine3.7 AS build
RUN foo
FROM build as unit-test
RUN bar
FROM alpine:3.7
RUN baz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"no untagged",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"no untagged with name",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM debian AS builder`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"other untagged cases are not ok",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM golang:1.9.3-alpine3.7 AS build
RUN foo
FROM node as unit-test
RUN bar
FROM alpine:3.7
RUN baz`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertContainsViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"scratch",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM scratch`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"untagged digest is not an error",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM ruby@sha256:f1dbca0f5dbc9`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"untagged digest is not an error",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM ruby:2`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3006")
		},
	)

	t.Run(
		"using args is not an error",
		func(t *testing.T) {
			t.Parallel()

			dockerfile := `FROM ${VALUE}`
			violations := testutils.LintDockerfile(dockerfile, allRules)

			testutils.AssertNoViolation(t, violations, "DL3006")
		},
	)
}
