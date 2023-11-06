package gitlet

import (
	"fmt"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

func Push(cCtx *cli.Context) error {
	stdOut, stdErr := RunCommand("git push", settings.Path)
	fmt.Println(stdOut)
	fmt.Println(stdErr)
	return nil
}
