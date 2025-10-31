package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/urfave/cli/v3"

	"github.com/farcloser/godolint/internal/parser"
	"github.com/farcloser/godolint/internal/process"
	"github.com/farcloser/godolint/internal/rule"
	"github.com/farcloser/godolint/internal/rules"
)

func main() {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	cmd := &cli.Command{
		Name:  "godolint",
		Usage: "Dockerfile linter",
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
			processor := process.NewProcessor([]rule.Rule{
				rules.DL4000(),
			})

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
