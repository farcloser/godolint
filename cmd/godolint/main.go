// Package main implements the godolint CLI for linting Dockerfiles.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"time"

	"github.com/alecthomas/kong"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/shell"
	"github.com/farcloser/godolint/sdk"
)

// errViolations is how Run tells main that the lint found something: the
// failures are already on stdout, and the exit status is 1 without a log line.
var errViolations = errors.New("violations found")

// CLI is the kong grammar: the flags and the Dockerfile arguments. Paths stay
// exactly as given — no kong file type, whose tilde and absolute-path
// expansion would rewrite the File field of every reported failure.
type CLI struct {
	DisableIgnorePragma bool     `help:"Disable inline ignore pragmas (# hadolint ignore=DLxxxx)."`
	WithoutShellcheck   bool     `help:"Disable shellcheck integration for RUN instruction validation."`
	Ignore              []string `help:"Rule code to ignore (repeatable: --ignore DL3006 --ignore SC2050)."                                                placeholder:"CODE"`
	ShellcheckRcfile    string   `help:"Shellcheckrc forwarded to shellcheck (--rcfile) when validating RUN instructions (requires shellcheck >= 0.10.0)." placeholder:"FILE"`

	Dockerfiles []string `arg:"" help:"Dockerfile(s) to lint." name:"dockerfile"`
}

// Run lints every Dockerfile, writes the failures to stdout as JSON, and
// returns errViolations when there were any.
func (c *CLI) Run() error {
	rules, err := c.buildRules()
	if err != nil {
		return err
	}

	// One processor, every rule, reused across files.
	processor := process.NewProcessor(rules).WithDisableIgnorePragmas(c.DisableIgnorePragma)

	allFailures, err := lintFiles(processor, c.Dockerfiles)
	if err != nil {
		return err
	}

	allFailures = dropIgnored(allFailures, c.Ignore)

	if err := json.NewEncoder(os.Stdout).Encode(allFailures); err != nil {
		return fmt.Errorf("failed to encode failures: %w", err)
	}

	if len(allFailures) > 0 {
		return errViolations
	}

	return nil
}

// buildRules assembles the rule set, wiring in the shellcheck integration
// unless it is disabled or the binary is missing from PATH.
func (c *CLI) buildRules() ([]rule.Rule, error) {
	rules := sdk.AllRules()

	if c.WithoutShellcheck {
		return rules, nil
	}

	// Missing shellcheck degrades gracefully (matching hadolint): warn and
	// lint without the integration rather than failing the run.
	//nolint:nilerr // intentional, see above.
	if _, err := exec.LookPath("shellcheck"); err != nil {
		log.Warn().Msg("shellcheck binary not found in PATH, shellcheck integration disabled")

		return rules, nil
	}

	checker := shell.NewBinaryShellchecker()
	// Fail fast on an unreadable rcfile: shellcheck errors are non-fatal per
	// rule (matching hadolint), so a bad path would otherwise silently
	// disable every SC check.
	if c.ShellcheckRcfile != "" {
		if _, err := os.Stat(c.ShellcheckRcfile); err != nil {
			return nil, fmt.Errorf("cannot read shellcheck rcfile: %w", err)
		}

		checker.RCFile = c.ShellcheckRcfile
	}

	return append(rules, shell.NewShellcheckRule(checker)), nil
}

// lintFiles runs the processor over each Dockerfile and returns the collected
// failures, each tagged with the file it came from.
func lintFiles(processor *process.Processor, paths []string) ([]rule.CheckFailure, error) {
	// Non-nil so an all-clean run still encodes as JSON [] rather than null.
	allFailures := []rule.CheckFailure{}

	for _, dockerfilePath := range paths {
		// #nosec G304 -- reading user-supplied Dockerfile paths is this tool's purpose.
		dockerfileContent, err := os.ReadFile(dockerfilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", dockerfilePath, err)
		}

		instructions, err := parser.NewBuildkitParser().Parse(dockerfileContent)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", dockerfilePath, err)
		}

		log.Debug().Str("file", dockerfilePath).Int("instructions", len(instructions)).Msg("Parsed Dockerfile")

		failures := processor.Run(instructions)
		for i := range failures {
			failures[i].File = dockerfilePath
		}

		allFailures = append(allFailures, failures...)
	}

	return allFailures, nil
}

// dropIgnored filters out failures whose rule code was --ignore'd.
func dropIgnored(failures []rule.CheckFailure, ignoredRules []string) []rule.CheckFailure {
	if len(ignoredRules) == 0 {
		return failures
	}

	filtered := []rule.CheckFailure{}

	for _, failure := range failures {
		if !slices.Contains(ignoredRules, string(failure.Code)) {
			filtered = append(filtered, failure)
		}
	}

	return filtered
}

func configureLogger(ctx context.Context, level ...zerolog.Level) {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger.WithContext(ctx)

	if len(level) > 0 {
		// Explicit level provided
		zerolog.SetGlobalLevel(level[0])
	} else {
		// Read from LOG_LEVEL environment variable
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}

		parsedLevel, err := zerolog.ParseLevel(logLevel)
		if err != nil {
			// Invalid level, default to info
			parsedLevel = zerolog.InfoLevel

			log.Warn().Str("LOG_LEVEL", logLevel).Msg("Invalid log level, defaulting to info")
		}

		zerolog.SetGlobalLevel(parsedLevel)
	}
}

func main() {
	configureLogger(context.Background())

	var cli CLI

	// Parse and usage errors are kong's: printed with the usage. Every
	// failure exits 1 — kong's own usage-error status (80) is folded into
	// the one non-zero status this tool has always had; --help stays 0.
	kctx := kong.Parse(&cli,
		kong.Name("godolint"),
		kong.Description("Dockerfile linter"),
		kong.UsageOnError(),
		kong.Exit(func(code int) {
			if code != 0 {
				code = 1
			}

			os.Exit(code)
		}),
	)

	if err := kctx.Run(); err != nil {
		if errors.Is(err, errViolations) {
			os.Exit(1)
		}

		log.Error().Err(err).Msg("failed to run godolint")
		os.Exit(1)
	}
}
