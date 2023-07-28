package injectzapr

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	_ "embed"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/draganm/inject/helpers"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

//go:embed template.go
var templateSource []byte

var Command = &cli.Command{
	Name:        "zapr",
	Description: "Inserts zapr v2 into main",
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

		mf.Body.List = append(filterOutCli(templateBlock.Body.List), mf.Body.List...)

		// findCliStatements(templateBlock.Body.List)

		tbf, err := findCliActionFunction(templateBlock.Body.List)
		if err != nil {
			return fmt.Errorf("could not find CLI action function in template: %w", err)
		}

		mbf, err := findCliActionFunction(mf.Body.List)

		if err != nil {
			return fmt.Errorf("could not find CLI action function in main: %w", err)
		}

		mbf.Body.List = append(tbf.Body.List, mbf.Body.List...)

		f.Imports = append(f.Imports, tf.Imports...)

		bb := &bytes.Buffer{}

		err = decorator.Fprint(bb, tf)

		if err != nil {
			return fmt.Errorf("could not serialize source: %w", err)
		}

		return os.WriteFile("main.go", bb.Bytes(), 0700)

	},
}

func isIdent(e dst.Expr, name string) bool {
	i, isIdent := e.(*dst.Ident)
	if !isIdent {
		return false
	}
	return i.Name == name

}

func isCallExpr(e dst.Expr, name string) bool {
	ce, isCallExpr := e.(*dst.CallExpr)
	if !isCallExpr {
		return false
	}

	se, isSelectorExpression := ce.Fun.(*dst.SelectorExpr)
	if !isSelectorExpression {
		return false
	}
	return se.Sel.Name == name
}

func filterOutCli(s []dst.Stmt) []dst.Stmt {
	return lo.Filter(s, func(st dst.Stmt, _ int) bool {
		fmt.Printf("%T\n", st)
		switch stmt := st.(type) {
		case *dst.AssignStmt:
			return !(len(stmt.Lhs) == 1 && isIdent(stmt.Lhs[0], "app"))
		case *dst.ExprStmt:
			return !isCallExpr(stmt.X, "RunAndExitOnError")
		}
		return true
	})
}

func findCliActionFunction(s []dst.Stmt) (*dst.FuncLit, error) {
	cliCall, found := lo.Find(s, func(st dst.Stmt) bool {
		switch stmt := st.(type) {
		case *dst.AssignStmt:
			return (len(stmt.Lhs) == 1 && isIdent(stmt.Lhs[0], "app"))
		}
		return false
	})

	if !found {
		return nil, errors.New("could not find cli assignment in template")
	}

	as, _ := cliCall.(*dst.AssignStmt)

	if len(as.Rhs) != 1 {
		return nil, errors.New("cli assignment in template is malformed")
	}

	fl, found := helpers.FindFirst[*dst.FuncLit](as)
	if !found {
		return nil, errors.New("cli assignment in template is malformed")
	}
	return fl, nil
}
