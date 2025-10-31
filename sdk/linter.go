package sdk

import (
	"context"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
)

// Linter performs Dockerfile linting.
type Linter struct {
	parser parser.Parser
	rules  []rule.Rule
}

// Option configures a Linter.
type Option func(*Linter)

// WithParser sets a custom parser implementation.
// By default, uses the buildkit parser.
func WithParser(p parser.Parser) Option {
	return func(l *Linter) {
		l.parser = p
	}
}

// WithRules sets custom rules to use for linting.
// By default, uses AllRules().
func WithRules(rules []rule.Rule) Option {
	return func(l *Linter) {
		l.rules = rules
	}
}

// WithRuleSet uses a predefined rule set.
func WithRuleSet(set RuleSet) Option {
	return func(l *Linter) {
		l.rules = GetRuleSet(set)
	}
}

// WithDisabledRules returns an option that disables specific rules by their codes.
// Example: WithDisabledRules("DL3000", "DL3007").
func WithDisabledRules(codes ...string) Option {
	return func(l *Linter) {
		l.rules = FilterRules(l.rules, codes)
	}
}

// WithShellcheck enables shellcheck integration for RUN instruction validation.
// Requires the shellcheck binary to be available in PATH.
func WithShellcheck() Option {
	return func(l *Linter) {
		checker := shell.NewBinaryShellchecker()
		scRule := shell.NewShellcheckRule(checker)
		l.rules = append(l.rules, scRule)
	}
}

// New creates a new Linter with the given options.
// By default, uses all implemented rules and the buildkit parser.
func New(opts ...Option) *Linter {
	lint := &Linter{
		parser: parser.NewBuildkitParser(),
		rules:  AllRules(),
	}

	for _, opt := range opts {
		opt(lint)
	}

	return lint
}

// Lint lints the given Dockerfile content.
// The context can be used for cancellation.
func (l *Linter) Lint(ctx context.Context, dockerfile []byte) (*Result, error) {
	// Check context cancellation before parsing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Parse Dockerfile
	instructions, err := l.parser.Parse(dockerfile)
	if err != nil {
		return nil, &ParseError{Err: err}
	}

	// Check context cancellation before processing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Run rules
	processor := process.NewProcessor(l.rules)
	failures := processor.Run(instructions)

	// Convert to SDK violations
	violations := make([]Violation, len(failures))
	for i, f := range failures {
		violations[i] = Violation{
			Code:     string(f.Code),
			Severity: convertSeverity(f.Severity),
			Message:  f.Message,
			Line:     f.Line,
		}
	}

	result := &Result{
		Violations: violations,
		Passed:     len(violations) == 0,
	}

	return result, nil
}

// LintFile is a convenience method that reads and lints a Dockerfile from a file path.
func (l *Linter) LintFile(ctx context.Context, path string) (*Result, error) {
	// Note: This would require os.ReadFile, but for now we keep the API surface simple
	// and let users read files themselves. This avoids adding file I/O concerns to the linter.
	// If we add this, we should also handle context cancellation properly.
	panic("LintFile not implemented - use Lint() with os.ReadFile()")
}

func convertSeverity(s rule.Severity) Severity {
	switch s {
	case rule.Error:
		return SeverityError
	case rule.Warning:
		return SeverityWarning
	case rule.Info:
		return SeverityInfo
	case rule.Style:
		return SeverityStyle
	default:
		return SeverityInfo
	}
}
