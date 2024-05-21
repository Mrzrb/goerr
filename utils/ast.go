package utils

import (
	"fmt"
	"go/ast"
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
