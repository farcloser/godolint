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

// Ported from CheckFailure in Hadolint/Rule.hs.
type CheckFailure struct {
	Code     RuleCode
	Severity Severity
	Message  string
	Line     int
}

// Ported from the concept of Rule in Hadolint/Rule.hs.
type Rule interface {
	// Code returns the unique rule identifier
	Code() RuleCode

	// Severity returns the severity level of violations
	Severity() Severity

	// Message returns the human-readable description of the rule
	Message() string

	// Check examines an instruction and returns whether it passes the rule.
	// Returns true if the instruction passes, false if it violates the rule.
	// Ported from the check function in Haskell simpleRule
	Check(instruction syntax.Instruction) bool
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

func (r *SimpleRule) Check(instruction syntax.Instruction) bool {
	return r.checker(instruction)
}
