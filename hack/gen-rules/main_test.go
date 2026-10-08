package main

import (
	"os"
	"path/filepath"
	"testing"
)

// A Haskell string gap (backslash, line break and indentation, backslash)
// is deleted whole: the message keeps the one space before the gap, not the
// indentation inside it.
func TestParseRuleFileCollapsesStringGaps(t *testing.T) {
	t.Parallel()

	src := `rule :: Rule args
rule = simpleRule code severity message check
  where
    code = "DL3007"
    severity = DLWarningC
    message =
      "Using latest is prone to errors if the image will ever update. Pin the version explicitly \
      \to a release tag"

    check (From BaseImage {tag = Just t}) = t /= "latest"
    check _ = True
`

	path := filepath.Join(t.TempDir(), "DL3007.hs")
	if err := os.WriteFile(path, []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}

	rule, err := parseRuleFile(path)
	if err != nil {
		t.Fatal(err)
	}

	want := "Using latest is prone to errors if the image will ever update. Pin the version explicitly to a release tag"
	if rule.Message != want {
		t.Errorf("Message = %q\nwant      %q", rule.Message, want)
	}
}
