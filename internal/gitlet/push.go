package gitlet

import (
	"fmt"
	"strings"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

func Push(_ *cli.Context) error {
	stdOut, stdErr := RunCommand("git fetch", "")
	printOutput(stdOut, stdErr)

	stdOut, stdErr = RunCommand("git push", settings.Path)
	printOutput(stdOut, stdErr)
	return nil
}

func printOutput(stdOut string, stdErr error) {
	if stdErr != nil {
		fmt.Println(stdErr)
	}

	if strings.TrimSpace(stdOut) != "" {
		fmt.Println(stdOut)
	}
}
