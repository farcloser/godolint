package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// dl3026State tracks stage aliases to allow referencing previous stages.
type dl3026State struct {
	aliases map[string]bool // stage aliases that can be referenced
}

// DL3026Rule checks for allowed registries.
type DL3026Rule struct {
	allowedRegistries []string // Empty means all registries allowed
}

// TODO: Add configuration support to specify allowed registries.
func DL3026() rule.Rule {
	return &DL3026Rule{
		allowedRegistries: []string{}, // Empty = all allowed (for now)
	}
}

func (r *DL3026Rule) Code() rule.RuleCode {
	return DL3026Meta.Code
}

func (r *DL3026Rule) Severity() rule.Severity {
	return DL3026Meta.Severity
}

func (r *DL3026Rule) Message() string {
	return DL3026Meta.Message
}

func (r *DL3026Rule) InitialState() rule.State {
	return rule.EmptyState(dl3026State{
		aliases: make(map[string]bool),
	})
}

func (r *DL3026Rule) Check(line int, state rule.State, instruction syntax.Instruction) rule.State {
	s := state.Data.(dl3026State)

	from, ok := instruction.(*syntax.From)
	if !ok {
		return state
	}

	// Track alias
	if from.Image.Alias != nil {
		s.aliases[*from.Image.Alias] = true
	}

	// Check if image is allowed
	imageName := from.Image.Image
	registry := extractRegistry(imageName)

	// Check if this is a reference to a previous stage
	if s.aliases[imageName] {
		return state.ReplaceData(s)
	}

	// If no registries are configured, all are allowed
	if len(r.allowedRegistries) == 0 {
		return state.ReplaceData(s)
	}

	// Special case: scratch is always allowed
	if imageName == "scratch" {
		return state.ReplaceData(s)
	}

	// Check if registry is in allowlist
	if !r.isRegistryAllowed(registry) {
		return state.AddFailure(rule.CheckFailure{
			Code:     DL3026Meta.Code,
			Severity: DL3026Meta.Severity,
			Message:  DL3026Meta.Message,
			Line:     line,
		})
	}

	return state.ReplaceData(s)
}

func (r *DL3026Rule) Finalize(state rule.State) rule.State {
	return state
}

// extractRegistry extracts the registry from an image name.
// If no registry is specified, defaults to docker.io.
func extractRegistry(imageName string) string {
	// If image contains /, check if first part is a registry
	if idx := strings.Index(imageName, "/"); idx != -1 {
		firstPart := imageName[:idx]
		// If first part contains a . or :, it's likely a registry
		if strings.Contains(firstPart, ".") || strings.Contains(firstPart, ":") {
			return firstPart
		}
	}

	// No explicit registry, default to docker.io
	return "docker.io"
}

// - Full wildcard: "*".
func (r *DL3026Rule) isRegistryAllowed(registry string) bool {
	for _, allowed := range r.allowedRegistries {
		if matchRegistry(allowed, registry) {
			return true
		}
	}

	return false
}

// matchRegistry checks if a registry matches an allowed pattern.
func matchRegistry(allowed, registry string) bool {
	// Full wildcard
	if allowed == "*" {
		return true
	}

	// Wildcard suffix: *.example.com
	if strings.HasPrefix(allowed, "*.") {
		suffix := strings.TrimPrefix(allowed, "*")

		return strings.HasSuffix(registry, suffix)
	}

	// Wildcard prefix: example.*
	if strings.HasSuffix(allowed, ".*") {
		prefix := strings.TrimSuffix(allowed, ".*")

		return strings.HasPrefix(registry, prefix)
	}

	// Exact match
	return registry == allowed
}
