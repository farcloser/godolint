// Package main demonstrates how quark would integrate godolint SDK
// instead of shelling out to hadolint binary.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/farcloser/godolint/sdk"
)

// This example mirrors quark's audit.Auditor.AuditDockerfile implementation,
// showing how to replace exec.Command("hadolint") with direct SDK usage.

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <Dockerfile>\n", os.Args[0])
		os.Exit(1)
	}

	dockerfilePath := os.Args[1]

	// Read Dockerfile
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read Dockerfile: %v\n", err)
		os.Exit(1)
	}

	// Create linter with all rules (matching hadolint default behavior)
	linter := sdk.New()

	// Lint with context support (for cancellation)
	ctx := context.Background()

	result, err := linter.Lint(ctx, content)
	if err != nil {
		// Check for parse errors
		if parseErr, ok := err.(*sdk.ParseError); ok {
			fmt.Fprintf(os.Stderr, "Failed to parse Dockerfile: %v\n", parseErr)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "Linting failed: %v\n", err)
		os.Exit(1)
	}

	// Format results as JSON (matching hadolint --format json)
	violations := make([]map[string]any, len(result.Violations))
	for i, v := range result.Violations {
		violations[i] = map[string]any{
			"code":    v.Code,
			"message": v.Message,
			"line":    v.Line,
			"level":   string(v.Severity),
		}
	}

	// Output JSON array (hadolint-compatible format)
	output, _ := json.MarshalIndent(violations, "", "  ")
	fmt.Println(string(output))

	// Exit with non-zero if violations found (matching hadolint behavior)
	if !result.Passed {
		os.Exit(1)
	}
}
