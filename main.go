package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"

	"github.com/kaatinga/commit/internal/commit"
	"github.com/kaatinga/commit/internal/gitlet"
	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

var version = "unknown"

func init() {
	file, err := os.Open("VERSION")
	if err != nil {
		return
	}

	var data []byte
	data, err = io.ReadAll(file)
	if err != nil {
		return
	}

	if len(data) > 0 {
		version = string(bytes.TrimSpace(data))
	}
}

func main() {
	settings.Init()
	gitlet.Init()

	app := &cli.App{
		Name:           "A git commit CLI tool",
		Description:    "Commit helps to generate commit messages.",
		DefaultCommand: "commit",
		Version:        version,
		Compiled:       time.Now(),
		Authors: []*cli.Author{
			{
				Name: "Michael Gunkoff",
			},
		},
		HelpName: "commit",
		Usage:    "automatic commit message generator",
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
				Name:  "key",
				Usage: "provide a valid key to work with openAI API",
				Action: func(context *cli.Context, s string) error {
					if len(s) != 51 {
						return cli.Exit("invalid key ", 1)
					}
					settings.APIKey = s
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "provide a valid path to work with git repository",
				Action: func(context *cli.Context, s string) error {
					if s != "" {
						settings.Path = s
					}

					return settings.DefinePaths()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
