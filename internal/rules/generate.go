package rules

// Generate rule stubs and tests from hadolint Haskell sources.
//
// Run: go generate ./internal/rules
//
// This does two things:
// 1. Extracts metadata (code, severity, message) from hadolint .hs rule files
//    and generates Go stub files for unimplemented rules
// 2. Extracts test cases from hadolint test files and generates Go tests
//    for implemented rules
//
// Rule stubs include placeholder check functions that need manual implementation.
// Tests are automatically ported from hadolint's test suite.

//go:generate go run ../../tools/gen-rules/main.go ../../hadolint/src/Hadolint/Rule
//go:generate go run ../../tools/gen-tests/main.go ../../hadolint/test/Hadolint/Rule
