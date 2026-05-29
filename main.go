package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kaatinga/commit/internal/commit"
	"github.com/kaatinga/commit/internal/gitlet"
	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

var version = "dev"

func main() {
	app := &cli.App{
		Name:           "commit",
		Description:    "Commit helps to generate commit messages.",
		DefaultCommand: "commit",
		Compiled:       time.Now(),
		Version:        version,
		Authors: []*cli.Author{
			{
				Name: "Michael Gunkoff",
			},
		},
		HelpName: "commit",
		Usage:    "automatic commit message generator",
		// Flag destinations are bound during flag parsing, which runs before
		// Before, so settings.Path is already set when we open the repository.
		Before: func(*cli.Context) error {
			return gitlet.Open(settings.Path)
		},
		Commands: []*cli.Command{
			{
				Name:   "commit",
				Action: commit.Generate,
				Hidden: true,
			},
			{
				Name:   "push",
				Action: actionChain(commit.Generate, gitlet.Push),
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "key",
				Usage:       "Mistral API key (falls back to the MISTRAL_API_KEY env var)",
				EnvVars:     []string{"MISTRAL_API_KEY"},
				Destination: &settings.APIKey,
			},
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path to the git repository",
				Value:       ".",
				Destination: &settings.Path,
			},
			&cli.BoolFlag{
				Name:        "dry-run",
				Usage:       "print the generated commit message without committing",
				EnvVars:     []string{"COMMIT_DRYRUN"},
				Destination: &settings.DryRun,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
