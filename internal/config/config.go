package config

// LabelType defines the validation type for a label.
type LabelType string

const (
	LabelTypeEmail   LabelType = "email"
	LabelTypeGitHash LabelType = "git-hash"
	LabelTypeRawText LabelType = "raw-text"
	LabelTypeRFC3339 LabelType = "rfc3339"
	LabelTypeSemVer  LabelType = "semver"
	LabelTypeSPDX    LabelType = "spdx"
	LabelTypeURL     LabelType = "url"
)

// Config holds configuration for godolint.
type Config struct {
	// AllowedRegistries is a list of allowed Docker registries.
	// Empty list means all registries are allowed.
	// Supports wildcards: "*", "*.example.com", "example.*"
	AllowedRegistries []string

	// LabelSchema defines required/optional labels and their types.
	// Key is the label name, value is the validation type.
	LabelSchema map[string]LabelType

	// StrictLabels when true, only labels in the schema are allowed.
	// When false, any labels are allowed (but schema labels are still validated).
	StrictLabels bool
}

// Default returns a default configuration (all rules permissive).
func Default() *Config {
	return &Config{
		AllowedRegistries: []string{}, // Empty = all allowed
		LabelSchema:       make(map[string]LabelType),
		StrictLabels:      false,
	}
}
