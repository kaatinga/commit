package gitlet

import (
	"fmt"
	"github.com/go-git/go-git/v5"

	"github.com/urfave/cli/v2"
)

func Push(cCtx *cli.Context) error {
	err := Repo.Push(&git.PushOptions{})
	if err != nil {
		return fmt.Errorf("error pushing to remote: %w", err)
	}

	fmt.Println("🚀 Commit pushed!")
	return nil
}
