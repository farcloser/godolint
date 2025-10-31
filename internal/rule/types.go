// Package rule defines the core types for hadolint rules.
// Ported from Hadolint/Rule.hs
package rule

import "github.com/farcloser/godolint/internal/syntax"

// Ported from DLSeverity in Hadolint/Rule.hs.
type Severity int

const (
	Error Severity = iota
	Warning
	Info
	Style
	Ignore
)

// String returns the string representation of the severity.
func (s Severity) String() string {
	switch s {
	case Error:
		return "error"
	case Warning:
		return "warning"
	case Info:
		return "info"
	case Style:
		return "style"
	case Ignore:
		return "ignore"
	default:
		return "unknown"
	}
}

// Ported from RuleCode in Hadolint/Rule.hs.
type RuleCode string

// RuleMeta contains metadata extracted from hadolint rule definitions.
// Used by generated code to separate metadata from implementation.
type RuleMeta struct {
	Code     RuleCode
	Severity Severity
	Message  string
}

// Ported from CheckFailure in Hadolint/Rule.hs.
type CheckFailure struct {
	Code     RuleCode
	Severity Severity
	Message  string
	Line     int
}

// State holds failures and custom state for a rule.
// Ported from State a in Hadolint/Rule.hs.
type State struct {
	Failures []CheckFailure
	Data     interface{} // Custom state data
}

// EmptyState creates a new state with no failures and the given initial data.
func EmptyState(data interface{}) State {
	return State{
		Failures: nil,
		Data:     data,
	}
}

// AddFailure adds a failure to the state.
func (s State) AddFailure(failure CheckFailure) State {
	return State{
		Failures: append(s.Failures, failure),
		Data:     s.Data,
	}
}

// ReplaceData replaces the custom state data.
func (s State) ReplaceData(data interface{}) State {
	return State{
		Failures: s.Failures,
		Data:     data,
	}
}

// Ported from the concept of Rule in Hadolint/Rule.hs.
// All rules are stateful - simple rules just use empty state.
type Rule interface {
	// Code returns the unique rule identifier
	Code() RuleCode

	// Severity returns the severity level of violations
	Severity() Severity

	// Message returns the human-readable description of the rule
	Message() string

	// InitialState returns the initial state for this rule
	InitialState() State

	// Check examines an instruction with the current state and returns updated state.
	// Ported from the step function in customRule.
	Check(line int, state State, instruction syntax.Instruction) State

	// Finalize processes the final state after all instructions have been checked.
	// Returns the updated state with any final failures.
	// Ported from markFailures in veryCustomRule.
	// Most rules don't need this and can just return the state unchanged.
	Finalize(state State) State
}

// Ported from simpleRule in Hadolint/Rule.hs.
type SimpleRule struct {
	code     RuleCode
	severity Severity
	message  string
	checker  func(syntax.Instruction) bool
}

// NewSimpleRule creates a new simple rule.
func NewSimpleRule(
	code RuleCode,
	severity Severity,
	message string,
	checker func(syntax.Instruction) bool,
) *SimpleRule {
	return &SimpleRule{
		code:     code,
		severity: severity,
		message:  message,
		checker:  checker,
	}
}

func (r *SimpleRule) Code() RuleCode {
	return r.code
}

func (r *SimpleRule) Severity() Severity {
	return r.severity
}

func (r *SimpleRule) Message() string {
	return r.message
}

func (r *SimpleRule) InitialState() State {
	return EmptyState(nil)
}

func (r *SimpleRule) Check(line int, state State, instruction syntax.Instruction) State {
	if !r.checker(instruction) {
		// Checker failed, add failure to state
		return state.AddFailure(CheckFailure{
			Code:     r.code,
			Severity: r.severity,
			Message:  r.message,
			Line:     line,
		})
	}

	return state
}

func (r *SimpleRule) Finalize(state State) State {
	return state // Simple rules don't need finalization
}
