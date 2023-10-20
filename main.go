package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/kaatinga/commit/internal/commit"
	"github.com/kaatinga/commit/internal/settings"
)

func init() {

}

func main() {
	app := &cli.App{
		Name:           "A git commit CLI tool",
		Description:    "Commit helps to generate commit messages.",
		DefaultCommand: "commit",
		Version:        "1.2.0",
		Compiled:       time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Michael Gunkoff",
				Email: "kaatinga@gmail.com",
			},
		},
		Copyright: "(c) Michael Gunkoff",
		HelpName:  "commit",
		Usage:     "automatic commit message generator",
		Commands: []*cli.Command{
			{
				Name:   "commit",
				Action: commit.Generate,
				Hidden: true,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "key",
				Usage: "provide a valid key to work with chatGPT",
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
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
