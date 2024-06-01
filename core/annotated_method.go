package core

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Method struct {
	Node
	MethodIdent
}

// Call implements Callable.
func (fc *Method) Call(pkg string, receiver string, returns []string, params ...string) string {
	if receiver == "" {
		panic(fmt.Sprintf("call receiver can not be empty %+v", fc))
	}
	var b strings.Builder
	b.WriteString(strings.Join(returns, " ,"))
	b.WriteString(" = ")
	b.WriteString(receiver + ".")
	b.WriteString(fc.Name)
	b.WriteString("(")
	b.WriteString(strings.Join(params, " ,"))
	b.WriteString(")")

	return b.String()
}

// Id implements Identity.
func (fc *Method) Id() string {
	return fc.Meta().Dir() + fc.Name
}

// Nodes implements Annotated.
func (fc *Method) Nodes() annotation.Node {
	return fc.Node
}

func NewMethod(n annotation.Node) *Method {
	met := &Method{
		Node: Node{n},
		MethodIdent: MethodIdent{
			FuncIdent: FuncIdent{},
			Receiver:  Ident{},
		},
	}
	met.Name = met.extractMethod(n)
	met.Annotation = met.Annotations()

	met.WalkField(func(f *ast.Field) {
		p := Ident{
			AnnotationsMix: AnnotationsMix{
				Annotation: []annotation.Annotation{},
			},
			Name: "",
			Type: "",
		}
		p.Name, p.Type, p.Annotation = utils.ExtractField(n, f)
		p.Raw = f
		met.Param = append(met.Param, p)
	})
	met.WalkReturn(func(f *ast.Field) {
		p := Ident{
			AnnotationsMix: AnnotationsMix{
				Annotation: []annotation.Annotation{},
			},
			Name: "",
			Type: "",
		}
		p.Name, p.Type, p.Annotation = utils.ExtractField(n, f)
		p.Raw = f
		met.Retern = append(met.Retern, p)
	})

	met.WalkReceiver(func(f *ast.Field) {
		p := Ident{
			AnnotationsMix: AnnotationsMix{
				Annotation: []annotation.Annotation{},
			},
			Name: "",
			Type: "",
		}
		p.Name, p.Type, p.Annotation = utils.ExtractField(n, f)
		met.Receiver = p
	})

	return met
}

func (fc *Method) extractMethod(n annotation.Node) string {
	ft := utils.MustBool(annotation.CastNode[*ast.FuncDecl](n))
	return ft.Name.Name
}

func (s *Method) Imports() []string {
	im := []string{}
	s.WalkField(func(f *ast.Field) {
		im = append(im, s.Node.Import(f.Type)...)
	})
	s.WalkReturn(func(f *ast.Field) {
		im = append(im, s.Node.Import(f.Type)...)
	})
	return im
}

func (s *Method) WalkField(fn func(*ast.Field)) {
	for _, v := range s.ASTNode().(*ast.FuncDecl).Type.Params.List {
		fn(v)
	}
}

func (s *Method) WalkReturn(fn func(*ast.Field)) {
	results := s.ASTNode().(*ast.FuncDecl).Type.Results
	if results == nil {
		return
	}
	for _, v := range results.List {
		fn(v)
	}
}

func (s *Method) WalkReceiver(fn func(*ast.Field)) {
	recv := s.ASTNode().(*ast.FuncDecl).Recv
	if recv == nil {
		return
	}
	for _, v := range recv.List {
		fn(v)
	}
}

// Annotate implements Annotated.
func (m *Method) Annotate() []annotation.Annotation {
	return m.Annotations()
}

var _ Annotated = (*Method)(nil)
