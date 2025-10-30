// Package parser provides an abstracted interface for parsing Dockerfiles.
// This allows swapping parser implementations (moby/buildkit, asottile/dockerfile, etc.)
package parser

import "github.com/farcloser/godolint/internal/syntax"

// Parser defines the interface for parsing Dockerfiles into AST.
type Parser interface {
	// Parse takes a Dockerfile's contents and returns a list of instructions with line numbers.
	// Returns an error if the Dockerfile cannot be parsed.
	Parse(dockerfile []byte) ([]syntax.InstructionPos, error)
}
