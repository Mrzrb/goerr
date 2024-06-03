package core

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Func struct {
	Node
	FuncIdent
}

// Call implements Callable.
func (f *Func) Call(pkg string, receiver string, returns []string, params ...string) string {
	if receiver != "" {
		panic(fmt.Sprintf("call receiver can not be none empty %+v", f))
	}
	var b strings.Builder
	if len(returns) > 0 {
		b.WriteString(strings.Join(returns, " ,"))
		b.WriteString(" := ")
		if pkg != f.Meta().PackageName() {
			b.WriteString(f.Meta().PackageName() + ".")
		}
	}
	b.WriteString(f.Name)
	b.WriteString("(")
	b.WriteString(strings.Join(params, " ,"))
	b.WriteString(")")

	return b.String()
}

// Id implements Identity.
func (f *Func) Id() string {
	return f.Meta().PackageName() + "." + f.Name
}

// Nodes implements Annotated.
func (f *Func) Nodes() annotation.Node {
	return f.Node
}

// Annotate implements Annotated.
func (f *Func) Annotate() []annotation.Annotation {
	return f.Annotations()
}

var _ Annotated = (*Func)(nil)

func NewFunc(n annotation.Node) *Func {
	fc := &Func{
		Node: Node{n},
		FuncIdent: FuncIdent{
			AnnotationsMix: AnnotationsMix{},
			Name:           "",
			Param:          []Ident{},
			Retern:         []Ident{},
		},
	}

	fc.Name = fc.extractFunc(n)
	fc.Annotation = fc.Annotations()
	fc.WalkField(func(f *ast.Field) {
		p := Ident{
			AnnotationsMix: AnnotationsMix{
				Annotation: []annotation.Annotation{},
			},
			Name: "",
			Type: "",
		}
		p.Raw = f
		p.Name, p.Type, p.Annotation, p.IsPointer = fc.extractField(f)
		fc.Param = append(fc.Param, p)
	})

	fc.WalkReturn(func(f *ast.Field) {
		p := Ident{
			AnnotationsMix: AnnotationsMix{
				Annotation: []annotation.Annotation{},
			},
			Name: "",
			Type: "",
		}
		p.Name, p.Type, p.Annotation, p.IsPointer = fc.extractField(f)
		p.Raw = f
		fc.Retern = append(fc.Retern, p)
	})

	return fc
}

func (fc *Func) extractFunc(n annotation.Node) string {
	ft := utils.MustBool(annotation.CastNode[*ast.FuncDecl](n))
	return ft.Name.Name
}

func (s *Func) extractField(n *ast.Field) (string, string, []annotation.Annotation, bool) {
	return utils.ExtractFieldWithPointer(s.Node, n)
}

func (s *Func) WalkField(fn func(*ast.Field)) {
	for _, v := range s.ASTNode().(*ast.FuncDecl).Type.Params.List {
		fn(v)
	}
}

func (s *Func) WalkReturn(fn func(*ast.Field)) {
	results := s.ASTNode().(*ast.FuncDecl).Type.Results
	if results == nil {
		return
	}
	for _, v := range results.List {
		fn(v)
	}
}

func (s *Func) Imports() []string {
	im := []string{}
	s.WalkField(func(f *ast.Field) {
		im = append(im, s.Node.Import(f.Type)...)
	})
	s.WalkReturn(func(f *ast.Field) {
		im = append(im, s.Node.Import(f.Type)...)
	})
	return im
}
