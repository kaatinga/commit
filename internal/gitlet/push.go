package gitlet

import (
	"fmt"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

func Push(_ *cli.Context) error {
	stdOut, stdErr := RunCommand("git fetch", settings.Path)
	printOutput(stdOut)
	printOutput(stdErr)

	stdOut, stdErr = RunCommand("git push", settings.Path)
	printOutput(stdOut)
	printOutput(stdErr)
	return nil
}

func printOutput(stdOut interface{}) {
	if stdOut == nil {
		return
	}
	fmt.Println(stdOut)
}
