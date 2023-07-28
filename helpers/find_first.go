package helpers

import (
	"reflect"

	"github.com/dave/dst"
)

type findByType struct {
	t     reflect.Type
	found dst.Node
}

func (v *findByType) Visit(node dst.Node) (w dst.Visitor) {
	if node == nil {
		return v
	}
	nv := reflect.ValueOf(node)
	if v.t.AssignableTo(nv.Type()) {
		v.found = node
		return nil
	}
	return v
}

func FindFirst[T any](node dst.Node) (T, bool) {
	fbt := &findByType{
		t: reflect.ValueOf(new(T)).Elem().Type(),
	}

	dst.Walk(fbt, node)

	if fbt.found == nil {
		return *new(T), false
	}

	return fbt.found.(T), true
}
