package sdk

// This file re-exports internal types to enable advanced usage:
// - Custom rule implementation
// - Custom parser implementation
// - Direct AST manipulation
// - Shell script analysis

import (
	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/internal/syntax"
)

// --- Rule Interface and Types ---

// Rule is the interface that all linting rules must implement.
// Users can implement this interface to create custom rules.
//
// Example:
//
//	type CustomRule struct{}
//
//	func (r *CustomRule) Code() RuleCode { return "CUSTOM001" }
//	func (r *CustomRule) Severity() RuleSeverity { return RuleSeverityError }
//	func (r *CustomRule) Message() string { return "Custom rule message" }
//	func (r *CustomRule) InitialState() RuleState { return EmptyRuleState(nil) }
//	func (r *CustomRule) Check(line int, state RuleState, instr Instruction) RuleState {
//	    // Implementation
//	    return state
//	}
//	func (r *CustomRule) Finalize(state RuleState) RuleState { return state }
type Rule = rule.Rule

// RuleCode is a unique identifier for a rule (e.g., "DL3000").
type RuleCode = rule.RuleCode

// RuleSeverity represents the severity level of a rule.
type RuleSeverity = rule.Severity

const (
	// RuleSeverityError indicates a critical issue.
	RuleSeverityError = rule.Error
	// RuleSeverityWarning indicates a potential issue.
	RuleSeverityWarning = rule.Warning
	// RuleSeverityInfo indicates informational feedback.
	RuleSeverityInfo = rule.Info
	// RuleSeverityStyle indicates a style preference.
	RuleSeverityStyle = rule.Style
	// RuleSeverityIgnore indicates the rule should be ignored.
	RuleSeverityIgnore = rule.Ignore
)

// RuleState holds the state for a rule during execution.
type RuleState = rule.State

// RuleFailure represents a single rule violation.
type RuleFailure = rule.CheckFailure

// EmptyRuleState creates a new rule state with no failures.
func EmptyRuleState(data interface{}) RuleState {
	return rule.EmptyState(data)
}

// NewSimpleRule creates a simple stateless rule.
// The checker function should return true if the instruction passes,
// false if it violates the rule.
func NewSimpleRule(
	code RuleCode,
	severity RuleSeverity,
	message string,
	checker func(Instruction) bool,
) Rule {
	return rule.NewSimpleRule(code, severity, message, checker)
}

// --- Parser Interface and Types ---

// Parser is the interface for Dockerfile parsers.
type Parser = parser.Parser

// NewBuildkitParser creates a parser using moby/buildkit.
func NewBuildkitParser() Parser {
	return parser.NewBuildkitParser()
}

// --- AST Types ---

// Instruction is the base interface for all Dockerfile instructions.
type Instruction = syntax.Instruction

// InstructionPos pairs an instruction with its line number.
type InstructionPos = syntax.InstructionPos

// Concrete instruction types for type assertions in custom rules.
type (
	FromInstruction        = syntax.From
	RunInstruction         = syntax.Run
	CopyInstruction        = syntax.Copy
	AddInstruction         = syntax.Add
	EnvInstruction         = syntax.Env
	LabelInstruction       = syntax.Label
	WorkdirInstruction     = syntax.Workdir
	UserInstruction        = syntax.User
	ExposeInstruction      = syntax.Expose
	VolumeInstruction      = syntax.Volume
	CmdInstruction         = syntax.Cmd
	EntrypointInstruction  = syntax.Entrypoint
	HealthcheckInstruction = syntax.Healthcheck
	MaintainerInstruction  = syntax.Maintainer
	ArgInstruction         = syntax.Arg
	StopSignalInstruction  = syntax.StopSignal
	ShellInstruction       = syntax.Shell
	OnBuildInstruction     = syntax.OnBuild
)

// Supporting types
type (
	BaseImage = syntax.BaseImage
	EnvPair   = syntax.EnvPair
	LabelPair = syntax.LabelPair
)

// --- Shell Analysis ---

// ShellOpts configures shell script analysis.
type ShellOpts = shell.ShellOpts

// DefaultShellOpts returns default shell options with common proxy variables.
func DefaultShellOpts() ShellOpts {
	return shell.DefaultShellOpts()
}

// ParsedShell represents a parsed shell script.
type ParsedShell = shell.ParsedShell

// Command represents a single command in a shell script.
type Command = shell.Command

// CmdPart represents a command argument or flag.
type CmdPart = shell.CmdPart

// ParseShell parses a shell script into commands.
func ParseShell(script string) (*ParsedShell, error) {
	return shell.ParseShell(script)
}

// Shell parsing helper functions
var (
	// FindCommandNames extracts all command names from parsed shell.
	FindCommandNames = shell.FindCommandNames
	// HasFlag checks if a command has a specific flag.
	HasFlag = shell.HasFlag
	// GetArgs returns all arguments from a command.
	GetArgs = shell.GetArgs
	// GetArgsNoFlags returns arguments excluding flags.
	GetArgsNoFlags = shell.GetArgsNoFlags
	// CountCommands returns the number of commands in parsed shell.
	CountCommands = shell.CountCommands
	// HasAnyFlag checks if command has any of the given flags.
	HasAnyFlag = shell.HasAnyFlag
	// GetFlagArg extracts the argument value for a flag.
	GetFlagArg = shell.GetFlagArg
)

// --- Shellcheck Integration ---

// Shellchecker is the interface for shellcheck implementations.
type Shellchecker = shell.Shellchecker

// NewBinaryShellchecker creates a shellchecker that uses the shellcheck binary.
func NewBinaryShellchecker() Shellchecker {
	return shell.NewBinaryShellchecker()
}

// NewNoopShellchecker creates a no-op shellchecker (for testing or when shellcheck unavailable).
func NewNoopShellchecker() Shellchecker {
	return shell.NewNoopShellchecker()
}

// NewShellcheckRule creates a shellcheck rule with the given checker.
func NewShellcheckRule(checker Shellchecker) Rule {
	return shell.NewShellcheckRule(checker)
}
