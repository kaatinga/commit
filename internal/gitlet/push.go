package gitlet

import (
	"github.com/kaatinga/commit/internal/settings"
	"github.com/urfave/cli/v2"
)

func Push(callback func(cCtx *cli.Context) error) func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {
		err := callback(cCtx)
		if err != nil {
			return err
		}

		_, err = RunCommand("git push", settings.Path)
		return err
	}
}
