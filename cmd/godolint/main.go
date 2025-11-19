// Package main implements the godolint CLI for linting Dockerfiles.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/sdk"
)

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
	ctx := context.Background()
	configureLogger(ctx)

	cmd := &cli.Command{
		Name:  "godolint",
		Usage: "Dockerfile linter",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "disable-ignore-pragma",
				Usage: "Disable inline ignore pragmas `# hadolint ignore=DLxxxx`",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errors.New("exactly one argument required: path to Dockerfile")
			}

			dockerfilePath := cmd.Args().First()

			// Read Dockerfile
			dockerfileContent, err := os.ReadFile(dockerfilePath)
			if err != nil {
				return fmt.Errorf("failed to read Dockerfile: %w", err)
			}

			// Parse Dockerfile using buildkit parser
			p := parser.NewBuildkitParser()
			instructions, err := p.Parse(dockerfileContent)
			if err != nil {
				return fmt.Errorf("failed to parse Dockerfile: %w", err)
			}

			log.Debug().Int("instructions", len(instructions)).Msg("Parsed Dockerfile")

			// Create processor with all rules
			processor := process.NewProcessor(sdk.AllRules()).
				WithDisableIgnorePragmas(cmd.Bool("disable-ignore-pragma"))

			// Run rules
			failures := processor.Run(instructions)

			// Output failures as JSON
			if err := json.NewEncoder(os.Stdout).Encode(failures); err != nil {
				return fmt.Errorf("failed to encode failures: %w", err)
			}

			// Exit with code 1 if any failures found
			if len(failures) > 0 {
				os.Exit(1)
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Error().Err(err).Msg("failed to run godolint")
		os.Exit(1)
	}
}
