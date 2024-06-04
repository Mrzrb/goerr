package core

import (
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Interface struct {
	Node
	Ident
}

// Annotate implements Annotated.
func (i *Interface) Annotate() []annotation.Annotation {
	panic("unimplemented")
}

// DstFileName implements Annotated.
// Subtle: this method shadows the method (Node).DstFileName of Interface.Node.
func (i *Interface) DstFileName(...string) string {
	panic("unimplemented")
}

// Nodes implements Annotated.
func (i *Interface) Nodes() annotation.Node {
	panic("unimplemented")
}

var _ Annotated = (*Interface)(nil)

func NewInterface(n annotation.Node) *Interface {
	node := &Interface{
		Node: Node{n},
		Ident: Ident{
			AnnotationsMix: AnnotationsMix{},
			Name:           "",
			Type:           "",
			Raw:            nil,
			IsPointer:      false,
			Package:        "",
		},
	}

	node.extractInterface(n)
	return node
}

func (i *Interface) extractInterface(node annotation.Node) {
	n := node.ASTNode().(*ast.TypeSpec)
	i.Name = n.Name.Name
	i.Type = n.Name.Name
	i.Package = node.Meta().PackageName()
	i.Raw = n
	i.Node.Node = node
}
