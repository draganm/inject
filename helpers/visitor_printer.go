package helpers

import (
	"fmt"

	"github.com/dave/dst"
)

type VisitorPrinter struct {
}

func (v VisitorPrinter) Visit(node dst.Node) (w dst.Visitor) {
	fmt.Printf("%T %v\n", node, node)
	return v
}
