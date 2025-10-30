// Package shell provides shellcheck integration for validating shell commands in RUN instructions.
package shell

import "github.com/farcloser/godolint/internal/rule"

// Shellchecker defines the interface for shellcheck integration.
// This allows rules to validate shell scripts in RUN instructions.
type Shellchecker interface {
	// Check runs shellcheck on a shell script and returns any violations found.
	Check(script string) ([]rule.CheckFailure, error)
}

// NoopShellchecker is a no-op implementation of Shellchecker.
// TODO: Implement actual shellcheck integration.
type NoopShellchecker struct{}

// NewNoopShellchecker creates a new no-op shellchecker.
func NewNoopShellchecker() *NoopShellchecker {
	return &NoopShellchecker{}
}

// Check always returns nil, indicating no violations found.
func (n *NoopShellchecker) Check(script string) ([]rule.CheckFailure, error) {
	return nil, nil
}
