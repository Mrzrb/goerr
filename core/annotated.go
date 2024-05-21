package core

import (
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Annotated interface {
	Generator
	Annotate() []annotation.Annotation
	Nodes() annotation.Node
}

type Generator interface {
	DstFileName() string
}

type AnnotationsMix struct {
	Annotation []annotation.Annotation
}

type Node struct {
	annotation.Node
}

func (n *Node) DstFileName() string {
	rawFileName := n.Meta().Dir() + "/" + n.Meta().FileName()
	fnames := strings.Split(rawFileName, ".")
	utils.InsertInPos(&fnames, "gen", len(fnames)-1)
	return strings.Join(fnames, ".")
}

type Ident struct {
	AnnotationsMix
	Name string
	Type string
	Raw  ast.Node
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
		if n.Recv != nil && len(n.Recv.List) > 0 {
			return NewMethod(node)
		}
		return NewFunc(node)
	}

	return nil
}

func ContainsAnnotate[T any](n Annotated) bool {
	return len(annotation.FindAnnotations[T](n.Annotate())) > 0
}
