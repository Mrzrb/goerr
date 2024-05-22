package core

import (
	"go/ast"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Method struct {
	Node
	MethodIdent
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
