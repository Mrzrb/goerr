package core

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Identity interface {
	Id() string
}

type Callable interface {
	Call(pkg string, receiver string, returns []string, params ...string) string
}

type Annotated interface {
	Generator
	Annotate() []annotation.Annotation
	Nodes() annotation.Node
}

type Generator interface {
	DstFileName(...string) string
}

type AnnotationsMix struct {
	Annotation []annotation.Annotation
}

type Node struct {
	annotation.Node
}

func (n *Node) Import(t ast.Expr) []string {
	i := utils.DistinctImports{}
	i.Merge(utils.GetImports(t, utils.ImprtTry(n)))
	return utils.Map(i.ToSlice(), func(t utils.Import) string {
		return t.Package
	})
}

func (n *Node) DstFileName(f ...string) string {
	rawFileName := n.Meta().Dir() + "/" + n.Meta().FileName()
	fnames := strings.Split(rawFileName, ".")
	utils.InsertInPos(&fnames, strings.Join(f, ".")+"gen", len(fnames)-1)
	return strings.Join(fnames, ".")
}

type Ident struct {
	AnnotationsMix
	Name      string
	Type      string
	Raw       ast.Node
	IsPointer bool
	Package   string
}

func (t *Ident) IsFromAnotherPkg() bool {
	return len(strings.Split(t.Type, ".")) > 1
}

func (t *Ident) PureType() string {
	return strings.Trim(t.Type, "*")
}

func (t *Ident) FullPackageType() string {
	if t.IsFromAnotherPkg() {
		return t.PureType()
	}

	return t.Package + "." + t.PureType()
}

func (i *Ident) Id(node annotation.Node) string {
	_, importPath, _ := findFieldPackage(i.Type, node)
	return fmt.Sprintf("%s.%s", importPath, strings.ReplaceAll(i.Type, "*", ""))
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

func Convert[T Annotated](node annotation.Node) (T, error) {
	var vv T
	n := Parse(node)
	if v, ok := Cast[T](n); ok {
		return v, nil
	}

	return vv, errors.New("cannot cast")
}

func ContainsAnnotate[T any](n Annotated) bool {
	return len(annotation.FindAnnotations[T](n.Annotate())) > 0
}

func GetByName(t []Annotated, name string) Annotated {
	a := utils.Filter(t, func(a Annotated) bool {
		n := ""
		switch a.(type) {
		case *Func:
			n = a.(*Func).Name
		case *Method:
			n = a.(*Method).Name
		case *Struct:
			n = a.(*Struct).Name
		}
		return n == name
	})
	if len(a) == 0 {
		panic(fmt.Sprintf("%s not found", name))
	}
	return a[0]
}

func GetMethod(t []Annotated, structName string) []Annotated {
	return utils.Filter(t, func(a Annotated) bool {
		n := ""
		switch a.(type) {
		case *Method:
			n = a.(*Method).Receiver.Type
		}
		return n == structName
	})
}
