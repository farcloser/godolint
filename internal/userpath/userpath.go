// Package userpath checks a path the user supplied before godolint uses it.
package userpath

import (
	"fmt"
	"path/filepath"

	"github.com/mycophonic/primordium/filesystem/pathcheck"
)

// Validate checks path with primordium's pathcheck: forbidden characters and
// the platform's length limits, component by component. pathcheck refuses a
// "." or ".." component, which a relative path the user types legitimately
// carries, so the path is checked in its absolute, cleaned form.
func Validate(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%q: %w", path, err)
	}

	if err = pathcheck.Validate(abs); err != nil {
		return fmt.Errorf("%q: %w", path, err)
	}

	return nil
}
