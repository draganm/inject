package injectcli

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/dave/dst/decorator"
	"github.com/draganm/inject/helpers"
	"github.com/urfave/cli/v2"
)

//go:embed template.go
var templateSource []byte

var Command = &cli.Command{
	Name:        "cli",
	Description: "Inserts urfave.cli v2 into main",
	Action: func(ctx *cli.Context) error {
		// fs := token.NewFileSet()
		mb, err := os.ReadFile("main.go")
		if err != nil {
			return fmt.Errorf("could not read main.go: %w", err)
		}

		f, err := decorator.Parse(mb)
		if err != nil {
			return fmt.Errorf("could not parse main.go: %w", err)
		}

		mf, err := helpers.FindMain(f)
		if err != nil {
			return fmt.Errorf("could not find main() in main.go: %w", err)
		}

		tf, err := decorator.Parse(templateSource)
		if err != nil {
			return fmt.Errorf("could not parse template file: %w", err)
		}

		templateBlock, err := helpers.FindTemplateFunction(tf)
		if err != nil {
			return err
		}

		mf.Body.List = append(templateBlock.Body.List, mf.Body.List...)
		f.Imports = append(f.Imports, tf.Imports...)
		// fmt.Println("imports", f.Imports)

		return decorator.Fprint(os.Stdout, tf)
		return nil
	},
}
