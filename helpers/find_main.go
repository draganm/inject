package helpers

import (
	"fmt"

	"github.com/dave/dst"
)

type findFunction struct {
	name                string
	functionDeclaration *dst.FuncDecl
}

func (v *findFunction) Visit(node dst.Node) (w dst.Visitor) {
	switch n := node.(type) {
	case *dst.FuncDecl:
		if n.Name.Name == v.name {
			v.functionDeclaration = n
		}
		return nil
	}
	return v
}

func FindMain(root dst.Node) (*dst.FuncDecl, error) {

	fmv := &findFunction{name: "main"}
	dst.Walk(fmv, root)
	if fmv.functionDeclaration != nil {
		return fmv.functionDeclaration, nil
	}

	return nil, fmt.Errorf("could not find main()")

}
