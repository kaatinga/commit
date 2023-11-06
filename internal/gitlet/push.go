package gitlet

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func Push(cCtx *cli.Context) error {
	err := Repo.Push(nil)
	if err != nil {
		return fmt.Errorf("error pushing to remote: %w", err)
	}

	fmt.Println("🚀 Commit pushed!")
	return nil
}
