package utils

import (
	"fmt"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

func ExtractTypeFromExpr(v ast.Expr) string {
	switch t := v.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return ExtractTypeFromExpr(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + ExtractTypeFromExpr(t.X)
	case *ast.ArrayType:
		return "[]" + ExtractTypeFromExpr(t.Elt)
	default:
		return fmt.Sprintf("%T", t)
	}
}

func ExtractField(s annotation.Node, n *ast.Field) (string, string, []annotation.Annotation) {
	annotatedNode := s.AnnotatedNode(n)
	var name, ty string
	if exp, ok := n.Type.(*ast.StarExpr); ok {
		name = exp.X.(*ast.Ident).Name
		ty = ExtractTypeFromExpr(n.Type)
	} else {
		if len(n.Names) > 0 {
			name = n.Names[0].Name
		}

		ty = ExtractTypeFromExpr(n.Type)
	}
	anns := annotatedNode.Annotations()
	return name, ty, anns
}
