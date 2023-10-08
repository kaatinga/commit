package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kaatinga/commit/internal/gitlet"
	"github.com/kaatinga/commit/internal/gpt"
	"github.com/urfave/cli/v2"
)

var (
	path   = "."
	apiKey = os.Getenv("OPENAI_API_KEY")
)

const (
	requestTemplate = `
Prepare a short commit message for the following changes, the response must contain only commit message itself. Code diff:
---------
`
)

func main() {
	app := &cli.App{
		Name:           "A git commit CLI tool",
		Description:    "Commit helps to generate commit messages.",
		DefaultCommand: "commit",
		Version:        "v19.99.0",
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
				Name: "commit",
				Action: func(cCtx *cli.Context) error {
					if apiKey == "" {
						return cli.Exit("openAI API key is not specified", 1)
					}

					diff, err := gitlet.GetDiff(path)
					if err != nil {
						return err
					}

					ctx, cancelFunc := context.WithTimeout(cCtx.Context, 10*time.Second)
					defer cancelFunc()

					var response string
					response, err = gpt.NewRequest(ctx, apiKey, requestTemplate+diff)
					if err != nil {
						return err
					}

					var gitInfo *gitlet.GitInfo
					gitInfo, err = gitlet.NewGitInfo(path, response)
					if err != nil {
						return err
					}

					fmt.Println("Added commit:\n", response)
					return gitInfo.Commit()
				},
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
					apiKey = s
					return nil
				},
			},
			&cli.StringFlag{
				Name:  "path",
				Usage: "provide a valid path to work with git repository",
				Action: func(context *cli.Context, s string) error {
					if s != "" {
						path = s
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
