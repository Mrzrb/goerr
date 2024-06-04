package core

import (
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Struct struct {
	Node
	Ident
	Field []Field
}

// Id implements Identity.
func (s *Struct) Id() string {
	return s.Meta().PackageName() + "." + s.Name
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

func (s *Struct) Annotate() []annotation.Annotation {
	return s.Annotations()
}

var _ Annotated = (*Struct)(nil)

func NewStruct(n annotation.Node) *Struct {
	node := &Struct{
		Node: Node{n},
		Ident: Ident{
			AnnotationsMix: AnnotationsMix{},
			Name:           "",
			Type:           "",
			Raw:            nil,
			Package:        n.Meta().PackageName(),
		},
		Field: []Field{},
	}
	// node ident
	node.Name, node.Type = node.extractStruct(n.ASTNode().(*ast.TypeSpec))
	node.Annotation = n.Annotations()
	node.Package = node.Meta().PackageName()
	node.WalkField(func(f *ast.Field) {
		nn, t, a, isPointer := node.extractField(f)
		fd := Field{
			Ident: Ident{
				AnnotationsMix: AnnotationsMix{Annotation: a},
				Name:           nn,
				Type:           t,
				Raw:            f,
				IsPointer:      isPointer,
				Package:        n.Meta().PackageName(),
			},
			Package:         "",
			Alias:           "",
			FullPackagePath: "",
			Parent:          n,
		}
		fd.Alias, fd.FullPackagePath, fd.Package = findFieldPackage(fd.Type, n)
		node.Field = append(node.Field, fd)
	})

	return node
}

func getPkgFromFullPath(path string) string {
	if path == "" {
		return ""
	}
	pathSlice := strings.Split(path, "/")
	if len(pathSlice) == 1 {
		return strings.Trim(pathSlice[0], "\"")
	}
	return strings.Trim(pathSlice[len(pathSlice)-1], "\"")
}

func findFieldPackage(tp string, n annotation.Node) (name string, importPath string, pkg string) {
	tp = strings.ReplaceAll(tp, "*", "")
	tpSlice := strings.Split(tp, ".")
	if len(tpSlice) <= 1 {
		return "", n.Meta().PackageName(), getPkgFromFullPath(n.Meta().PackageName())
	}

	utils.Walk(n.Imports(), func(is *ast.ImportSpec) {
		if is == nil {
			return
		}

		if is.Name != nil {
			if tpSlice[0] == is.Name.Name {
				name = tpSlice[0]
				importPath = is.Path.Value
				return
			}
		} else {
			path := strings.ReplaceAll(is.Path.Value, "\"", "")
			pathSlice := strings.Split(path, "/")
			if len(pathSlice) < 1 {
				return
			}
			if pathSlice[len(pathSlice)-1] == tpSlice[0] {
				name = tpSlice[0]
				importPath = is.Path.Value
				return
			}
		}
	})
	return name, strings.ReplaceAll(importPath, "\"", ""), getPkgFromFullPath(importPath)
}

func (s *Struct) extractStruct(n *ast.TypeSpec) (string, string) {
	return n.Name.Name, n.Name.Name
}

func (s *Struct) extractField(n *ast.Field) (string, string, []annotation.Annotation, bool) {
	return utils.ExtractFieldWithPointer(s.Node, n)
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

func (s *Struct) WalkFieldAnnoted(fn func(Field)) {
	for _, v := range s.Field {
		fn(v)
	}
}
