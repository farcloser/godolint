package userpath_test

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mycophonic/primordium/filesystem/pathcheck"

	"github.com/forkcloser/godolint/internal/userpath"
)

func TestValidateAccepts(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		"Dockerfile",
		"./Dockerfile",
		filepath.Join("..", "build", "Dockerfile"),
		filepath.Join(t.TempDir(), ".shellcheckrc"),
	} {
		if err := userpath.Validate(path); err != nil {
			t.Errorf("Validate(%q): %v", path, err)
		}
	}
}

func TestValidateRejects(t *testing.T) {
	t.Parallel()

	for _, path := range []string{
		"Docker\x00file",
		filepath.Join("dir", strings.Repeat("a", 256)),
	} {
		if err := userpath.Validate(path); !errors.Is(err, pathcheck.ErrInvalidPath) {
			t.Errorf("Validate(%q): %v, want pathcheck.ErrInvalidPath", path, err)
		}
	}
}
