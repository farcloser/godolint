package rules

import (
	"strings"

	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/syntax"
)

// Source: DL3020.hs.
func DL3020() rule.Rule {
	return rule.NewSimpleRule(
		"DL3020",
		rule.Error,
		"Use COPY instead of ADD for files and folders",
		checkDL3020,
	)
}

func checkDL3020(instruction syntax.Instruction) bool {
	add, ok := instruction.(*syntax.Add)
	if !ok {
		return true
	}

	// ADD should only be used for archives or URLs
	// For files and folders, use COPY instead
	for _, src := range add.Source {
		if !isArchive(src) && !isURL(src) {
			return false
		}
	}

	return true
}

// isArchive checks if a path has an archive file extension that Docker ADD auto-extracts.
// List matches hadolint's archiveFileFormatExtensions from Hadolint/Rule.hs
func isArchive(path string) bool {
	archiveExtensions := []string{
		".tar", ".Z", ".bz2", ".gz", ".lz", ".lzma",
		".tZ", ".tb2", ".tbz", ".tbz2", ".tgz",
		".tlz", ".tpz", ".txz", ".xz",
	}

	path = dropQuotes(path)
	for _, ext := range archiveExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

// isURL checks if a path is a URL.
func isURL(path string) bool {
	path = dropQuotes(path)

	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
