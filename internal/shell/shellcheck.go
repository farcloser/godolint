// Package shell provides shellcheck integration for validating shell commands in RUN instructions.
package shell

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// ShellOpts contains options for running shellcheck.
// Ported from Hadolint.Shell.ShellOpts.
type ShellOpts struct {
	ShellName string            // Shell command (e.g., "/bin/sh -c")
	EnvVars   map[string]string // Environment variables to export
}

// DefaultShellOpts returns the default shell options.
// Matches hadolint's defaultShellOpts with common proxy variables.
func DefaultShellOpts() ShellOpts {
	return ShellOpts{
		ShellName: "/bin/sh -c",
		EnvVars: map[string]string{
			"HTTP_PROXY":  "1",
			"http_proxy":  "1",
			"HTTPS_PROXY": "1",
			"https_proxy": "1",
			"FTP_PROXY":   "1",
			"ftp_proxy":   "1",
			"NO_PROXY":    "1",
			"no_proxy":    "1",
		},
	}
}

// Shellchecker defines the interface for shellcheck integration.
type Shellchecker interface {
	// Check runs shellcheck on a shell script with the given options.
	// Returns violations found by shellcheck.
	Check(script string, opts ShellOpts) ([]rule.CheckFailure, error)
}

// BinaryShellchecker shells out to the shellcheck binary.
type BinaryShellchecker struct{}

// NewBinaryShellchecker creates a shellchecker that uses the shellcheck binary.
func NewBinaryShellchecker() *BinaryShellchecker {
	return &BinaryShellchecker{}
}

// shellcheckOutput represents the JSON output from shellcheck.
type shellcheckOutput struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	EndLine int    `json:"endLine"`
	Column  int    `json:"column"`
	Level   string `json:"level"` // "error", "warning", "info", "style"
	Code    int    `json:"code"`  // SC code number
	Message string `json:"message"`
}

// Check runs shellcheck on the given script.
// Ported from Hadolint.Shell.shellcheck.
func (*BinaryShellchecker) Check(script string, opts ShellOpts) ([]rule.CheckFailure, error) {
	// Skip non-POSIX shells (pwsh, powershell, cmd)
	shellLower := strings.ToLower(opts.ShellName)
	if strings.Contains(shellLower, "pwsh") ||
		strings.Contains(shellLower, "powershell") ||
		strings.Contains(shellLower, "cmd") {
		return nil, nil
	}

	// Skip if script has unsupported shebang
	if hasUnsupportedShebang(script) {
		return nil, nil
	}

	// Build complete script with shebang and exports
	fullScript := buildScript(script, opts)

	// Write script to temp file
	tmpFile, err := os.CreateTemp("", "shellcheck-*.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(fullScript); err != nil {
		return nil, fmt.Errorf("failed to write script: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	// Run shellcheck with JSON output
	// Exclude codes like hadolint does:
	// - SC2187: ash shell not supported warning
	// - SC1090: can't follow sourced files (requires shell directives)
	// - SC1091: can't follow sourced files (requires shell directives)
	//nolint:gosec
	cmd := exec.Command("shellcheck",
		"--format=json",
		"--exclude=SC2187,SC1090,SC1091",
		"--severity=style", // Minimum severity (matches hadolint)
		tmpFile.Name())

	output, err := cmd.CombinedOutput()
	// shellcheck returns non-zero if violations found, which is expected
	// Only return error if we couldn't run shellcheck at all
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, fmt.Errorf("failed to run shellcheck: %w", err)
		}
	}

	// Parse JSON output
	var scResults []shellcheckOutput
	if len(output) > 0 {
		if err := json.Unmarshal(output, &scResults); err != nil {
			return nil, fmt.Errorf("failed to parse shellcheck output: %w", err)
		}
	}

	// Convert to CheckFailures
	var failures []rule.CheckFailure
	for _, sc := range scResults {
		failures = append(failures, rule.CheckFailure{
			Code:     rule.RuleCode(fmt.Sprintf("SC%d", sc.Code)),
			Severity: convertSeverity(sc.Level),
			Message:  sc.Message,
			Line:     0, // Will be set by the rule when it adds the failure
			Column:   1, // Hardcoded to 1 (matches hadolint)
		})
	}

	return failures, nil
}

// buildScript constructs the complete script to pass to shellcheck.
// Ported from the script construction in Hadolint.Shell.shellcheck.
func buildScript(runCommand string, opts ShellOpts) string {
	var build strings.Builder

	// Add shebang from shell option
	shebang := extractShell(opts.ShellName)
	if shebang == "" {
		shebang = "/bin/sh"
	}

	_, _ = build.WriteString("#!")
	_, _ = build.WriteString(shebang)
	_, _ = build.WriteString("\n")

	// Export environment variables
	for key := range opts.EnvVars {
		_, _ = build.WriteString("export ")
		_, _ = build.WriteString(key)
		_, _ = build.WriteString("=1\n")
	}

	// Add the actual RUN command
	_, _ = build.WriteString(runCommand)

	return build.String()
}

// extractShell extracts the shell path from shell command.
// Ported from extractShell in Hadolint.Shell.
func extractShell(shellCmd string) string {
	// Take first word from shell command
	// E.g., "/bin/bash -c" -> "/bin/bash"
	parts := strings.Fields(shellCmd)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

// hasUnsupportedShebang checks if script starts with an unsupported shebang.
// Ported from Hadolint.Shell.hasUnsupportedShebang.
func hasUnsupportedShebang(script string) bool {
	if !strings.HasPrefix(script, "#!") {
		return false
	}

	supported := []string{
		"#!/bin/sh",
		"#!/bin/bash",
		"#!/bin/ksh",
		"#!/usr/bin/env sh",
		"#!/usr/bin/env bash",
		"#!/usr/bin/env ksh",
	}

	for _, prefix := range supported {
		if strings.HasPrefix(script, prefix) {
			return false
		}
	}

	return true
}

// convertSeverity converts shellcheck severity to hadolint severity.
// Ported from Hadolint.Rule.Shellcheck.getDLSeverity.
func convertSeverity(scLevel string) rule.Severity {
	switch scLevel {
	case "warning":
		return rule.Warning
	case "info":
		return rule.Info
	case "style":
		return rule.Style
	default:
		return rule.Error
	}
}

// ShellcheckRule is a stateful rule that runs shellcheck on RUN instructions.
// Ported from Hadolint.Rule.Shellcheck.
type ShellcheckRule struct {
	checker Shellchecker
}

// shellState tracks shell options across instructions.
// Ported from Acc in Hadolint.Rule.Shellcheck.
type shellState struct {
	opts        ShellOpts
	defaultOpts ShellOpts
}

// NewShellcheckRule creates a new shellcheck rule.
func NewShellcheckRule(checker Shellchecker) *ShellcheckRule {
	return &ShellcheckRule{
		checker: checker,
	}
}

// Code returns the rule code.
func (*ShellcheckRule) Code() rule.RuleCode {
	return "SHELLCHECK"
}

// Severity returns the rule severity.
func (*ShellcheckRule) Severity() rule.Severity {
	// Shellcheck violations have their own severities
	return rule.Info
}

// Message returns the rule message.
func (*ShellcheckRule) Message() string {
	return "ShellCheck violations in RUN instructions"
}

// InitialState returns the initial state for this rule.
func (*ShellcheckRule) InitialState() rule.State {
	defaultOpts := DefaultShellOpts()

	return rule.EmptyState(shellState{
		opts:        defaultOpts,
		defaultOpts: defaultOpts,
	})
}

// Check processes each instruction and updates shell state.
// Ported from scrule in Hadolint.Rule.Shellcheck.
func (r *ShellcheckRule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	// Extract current shell state
	var shState shellState
	if state.Data != nil {
		shState = state.Data.(shellState)
	} else {
		defaultOpts := DefaultShellOpts()
		shState = shellState{
			opts:        defaultOpts,
			defaultOpts: defaultOpts,
		}
	}

	switch instr := instruction.(type) {
	case *syntax.From:
		// New stage - reset to default options
		return state.ReplaceData(shellState{
			opts:        shState.defaultOpts,
			defaultOpts: shState.defaultOpts,
		})

	case *syntax.Arg:
		// Add ARG to environment variables
		newOpts := shState.opts
		if newOpts.EnvVars == nil {
			newOpts.EnvVars = make(map[string]string)
		}
		// Copy existing vars
		envCopy := make(map[string]string)
		for k, v := range shState.opts.EnvVars {
			envCopy[k] = v
		}

		envCopy[instr.ArgName] = "1"
		newOpts.EnvVars = envCopy

		return state.ReplaceData(shellState{
			opts:        newOpts,
			defaultOpts: shState.defaultOpts,
		})

	case *syntax.Env:
		// Add ENV variables
		newOpts := shState.opts
		if newOpts.EnvVars == nil {
			newOpts.EnvVars = make(map[string]string)
		}
		// Copy existing vars
		envCopy := make(map[string]string)
		for k, v := range shState.opts.EnvVars {
			envCopy[k] = v
		}

		for _, pair := range instr.Pairs {
			envCopy[pair.Key] = "1"
		}

		newOpts.EnvVars = envCopy

		return state.ReplaceData(shellState{
			opts:        newOpts,
			defaultOpts: shState.defaultOpts,
		})

	case *syntax.Shell:
		// Update shell command
		if len(instr.Arguments) > 0 {
			shellCmd := strings.Join(instr.Arguments, " ")
			newOpts := shState.opts
			newOpts.ShellName = shellCmd

			return state.ReplaceData(shellState{
				opts:        newOpts,
				defaultOpts: shState.defaultOpts,
			})
		}

	case *syntax.Run:
		// Run shellcheck on the command
		violations, err := r.checker.Check(instr.Command, shState.opts)
		if err != nil {
			// Log error but don't fail the rule
			// (matching hadolint behavior - shellcheck failures are not fatal)
			return state
		}

		// Add all shellcheck violations to state with current line number
		newState := state

		for _, v := range violations {
			v.Line = line
			newState = newState.AddFailure(v)
		}

		return newState
	}

	return state
}

// Finalize performs final checks after processing all instructions.
func (*ShellcheckRule) Finalize(state rule.State) rule.State {
	return state // No finalization needed
}

// NoopShellchecker is a no-op implementation for when shellcheck is not available.
type NoopShellchecker struct{}

// NewNoopShellchecker creates a new no-op shellchecker.
func NewNoopShellchecker() *NoopShellchecker {
	return &NoopShellchecker{}
}

// Check always returns nil.
func (*NoopShellchecker) Check(script string, opts ShellOpts) ([]rule.CheckFailure, error) {
	return nil, nil
}
