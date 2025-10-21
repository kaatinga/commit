package gitlet

import (
	"fmt"

	"github.com/kaatinga/commit/internal/settings"

	"github.com/urfave/cli/v2"
)

func Push(_ *cli.Context) error {
	output, err := RunCommand("git fetch", "")
	if output != "" {
		fmt.Println(output)
	}
	if err != nil {
		return fmt.Errorf("git fetch failed: %w", err)
	}

	output, err = RunCommand("git push", settings.RepositoryPath)
	if output != "" {
		fmt.Println(output)
	}

	return err
}
