package main

import (
	"github.com/draganm/inject/command/injectcli"
	"github.com/draganm/inject/command/injectzapr"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Action: func(ctx *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			injectcli.Command,
			injectzapr.Command,
		},
	}
	app.RunAndExitOnError()
}
