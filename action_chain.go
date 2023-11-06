package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func actionChain(fns ...func(cCtx *cli.Context) error) func(cCtx *cli.Context) error {
	return func(cCtx *cli.Context) error {
		for _, fn := range fns {
			err := fn(cCtx)
			if err != nil {
				return fmt.Errorf("error executing a func in chain: %w", err)
			}
		}

		return nil
	}
}
