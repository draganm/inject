package injectcli

import (
	_ "embed"

	"github.com/urfave/cli/v2"
)

func template() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	app.RunAndExitOnError()
}
