package main

import (
	"github.com/draganm/inject/command/injectcli"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		// Action: func(ctx *cli.Context) error {

		// },
		Commands: []*cli.Command{
			injectcli.Command,
		},
	}
	app.RunAndExitOnError()
}
