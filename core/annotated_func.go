package core

import (
	"go/ast"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Func struct {
	annotation.Node
	FuncIdent
}

// Annotate implements Annotated.
func (f *Func) Annotate() []annotation.Annotation {
	return f.Annotations()
}

var _ Annotated = (*Func)(nil)

func NewFunc(n annotation.Node) *Func {
	fc := &Func{
		Node: n,
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
		p.Name, p.Type, p.Annotation = fc.extractField(f)
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
		p.Name, p.Type, p.Annotation = fc.extractField(f)
		fc.Retern = append(fc.Retern, p)
	})

	return fc
}

func (fc *Func) extractFunc(n annotation.Node) string {
	ft := utils.MustBool(annotation.CastNode[*ast.FuncDecl](n))
	return ft.Name.Name
}

func (s *Func) extractField(n *ast.Field) (string, string, []annotation.Annotation) {
	annotatedNode := s.AnnotatedNode(n)
	var name, ty string
	if exp, ok := n.Type.(*ast.StarExpr); ok {
		name = exp.X.(*ast.Ident).Name
		ty = utils.ExtractTypeFromExpr(n.Type)
	} else {
		if len(n.Names) > 0 {
			name = n.Names[0].Name
		}

		ty = utils.ExtractTypeFromExpr(n.Type)
	}
	anns := annotatedNode.Annotations()
	return name, ty, anns
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
