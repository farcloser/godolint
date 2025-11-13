// Package main demonstrates integration with shellcheck for RUN instruction validation.
package main

import (
	"fmt"
	"os"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
)

func main() {
	dockerfile := `FROM debian:bookworm
ENV MY_VAR=value
ARG BUILD_ARG=default
SHELL ["/bin/bash", "-c"]

# This should trigger SC2086 (unquoted variable)
RUN echo $MY_VAR

# This should trigger SC2154 (undefined variable)
RUN echo $UNDEFINED_VAR

# This should be fine
RUN apt-get update && apt-get install -y curl
`

	// Parse
	p := parser.NewBuildkitParser()

	instructions, err := p.Parse([]byte(dockerfile))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse: %v\n", err)
		os.Exit(1)
	}

	// Create shellcheck rule
	checker := shell.NewBinaryShellchecker()
	scRule := shell.NewShellcheckRule(checker)

	// Process
	processor := process.NewProcessor([]rule.Rule{scRule})
	failures := processor.Run(instructions)

	// Report
	fmt.Printf("Found %d shellcheck violations:\n\n", len(failures))

	for _, f := range failures {
		fmt.Printf("Line %d: [%s] %s - %s\n",
			f.Line, f.Severity, f.Code, f.Message)
	}
}
