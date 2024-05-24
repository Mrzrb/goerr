package core

import (
	"go/ast"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Struct struct {
	Node
	Ident
	Field []Field
}

func (s *Struct) Imports() []string {
	im := []string{}
	s.WalkField(func(f *ast.Field) {
		im = append(im, s.Node.Import(f.Type)...)
	})
	return im
}

// Nodes implements Annotated.
func (s *Struct) Nodes() annotation.Node {
	return s.Node
}

type Field struct {
	Ident
}

func (s *Struct) Annotate() []annotation.Annotation {
	return s.Annotations()
}

var _ Annotated = (*Struct)(nil)

func NewStruct(n annotation.Node) *Struct {
	node := &Struct{
		Node:  Node{n},
		Ident: Ident{},
		Field: []Field{},
	}
	// node ident
	node.Name, node.Type = node.extractStruct(n.ASTNode().(*ast.TypeSpec))
	node.Annotation = n.Annotations()
	node.WalkField(func(f *ast.Field) {
		n, t, a := node.extractField(f)
		fd := Field{
			Ident: Ident{
				AnnotationsMix: AnnotationsMix{Annotation: a},
				Name:           n,
				Type:           t,
				Raw:            f,
			},
		}
		node.Field = append(node.Field, fd)
	})

	return node
}

func (s *Struct) extractStruct(n *ast.TypeSpec) (string, string) {
	return n.Name.Name, n.Name.Name
}

func (s *Struct) extractField(n *ast.Field) (string, string, []annotation.Annotation) {
	return utils.ExtractField(s.Node, n)
	// annotatedNode := s.AnnotatedNode(n)
	// name, ty := n.Names[0].Name, utils.ExtractTypeFromExpr(n.Type)
	// anns := annotatedNode.Annotations()
	// return name, ty, anns
}

func (s *Struct) WalkField(fn func(*ast.Field)) {
	for _, v := range s.ASTNode().(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		fn(v)
	}
}
