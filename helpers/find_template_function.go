package helpers

import (
	"fmt"

	"github.com/dave/dst"
)

func FindTemplateFunction(f *dst.File) (*dst.FuncDecl, error) {

	fmv := &findFunction{name: "template"}
	dst.Walk(fmv, f)
	if fmv.functionDeclaration != nil {
		return fmv.functionDeclaration, nil
	}

	if fmv.functionDeclaration == nil {
		return nil, fmt.Errorf("could not find template() function")
	}

	return fmv.functionDeclaration, nil

}
