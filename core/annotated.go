package core

import (
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Annotated interface {
	Annotate() []annotation.Annotation
}

type AnnotationsMix struct {
	Annotation []annotation.Annotation
}

// @GetterSetter
type Ident struct {
	AnnotationsMix
	Name string
	Type string
}

type FuncIdent struct {
	AnnotationsMix
	Name   string
	Param  []Ident
	Retern []Ident
}

type MethodIdent struct {
	FuncIdent
	Receiver Ident
}

func Cast[T Annotated](n Annotated) (T, bool) {
	v, ok := n.(T)
	return v, ok
}

func Parse(node annotation.Node) Annotated {
	if _, ok := annotation.CastNode[*ast.TypeSpec](node); ok {
		return NewStruct(node)
	}
	if n, ok := annotation.CastNode[*ast.FuncDecl](node); ok {
		if len(n.Recv.List) > 0 {
			return NewMethod(node)
		}
		return NewFunc(node)
	}

	return nil
}
