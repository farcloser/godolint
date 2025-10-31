package sdk

// Severity represents the severity level of a violation.
type Severity string

const (
	// SeverityError indicates a critical issue that should be fixed.
	SeverityError Severity = "error"
	// SeverityWarning indicates a potential issue.
	SeverityWarning Severity = "warning"
	// SeverityInfo indicates informational feedback.
	SeverityInfo Severity = "info"
	// SeverityStyle indicates a style preference.
	SeverityStyle Severity = "style"
)

// Violation represents a single linting violation.
type Violation struct {
	// Code is the rule code (e.g., "DL3000", "SC2086").
	Code string `json:"code"`
	// Severity is the violation severity level.
	Severity Severity `json:"severity"`
	// Message is the human-readable description.
	Message string `json:"message"`
	// Line is the line number in the Dockerfile (1-indexed).
	Line int `json:"line"`
}

// Result contains the linting results.
type Result struct {
	// Violations contains all detected violations.
	Violations []Violation
	// Passed indicates whether the Dockerfile passed linting (no violations).
	Passed bool
}

// HasErrors returns true if any violations have Error severity.
func (r *Result) HasErrors() bool {
	for _, v := range r.Violations {
		if v.Severity == SeverityError {
			return true
		}
	}
	return false
}

// HasWarnings returns true if any violations have Warning severity.
func (r *Result) HasWarnings() bool {
	for _, v := range r.Violations {
		if v.Severity == SeverityWarning {
			return true
		}
	}
	return false
}

// CountBySeverity returns the count of violations for each severity level.
func (r *Result) CountBySeverity() map[Severity]int {
	counts := make(map[Severity]int)
	for _, v := range r.Violations {
		counts[v.Severity]++
	}
	return counts
}
