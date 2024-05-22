package utils

import (
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/utils/stream"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type ImportLookup func(alias string) (importPath string, found bool)

type Import struct {
	Alias   string
	Package string
}

type DistinctImports map[Import]struct{}

func NewDistinctImports() DistinctImports {
	return map[Import]struct{}{}
}

func (d DistinctImports) Append(i Import) {
	d[i] = struct{}{}
}

func (d DistinctImports) Merge(new DistinctImports) {
	for k := range new {
		d[k] = struct{}{}
	}
}

func (d DistinctImports) MergeSlice(new []Import) {
	for _, k := range new {
		d[k] = struct{}{}
	}
}

func (d DistinctImports) IsEmpty() bool {
	return len(d) == 0
}

func (d DistinctImports) ToSlice() []Import {
	if d.IsEmpty() {
		return nil
	}
	i := make([]Import, len(d))
	index := 0
	for k := range d {
		i[index] = k
		index++
	}
	return i
}

func getSelectorExpr(e ast.Node) *ast.SelectorExpr {
	switch i := e.(type) {
	case *ast.ArrayType:
		if in, ok := i.Elt.(*ast.StarExpr); ok {
			if inn, ok := in.X.(*ast.SelectorExpr); ok {
				return inn
			}
		}
	}

	return nil
}

func ImprtTry(node annotation.Node) ImportLookup {
	return func(alias string) (importPath string, found bool) {
		imp := stream.
			OfSlice(Map(node.Imports(), func(t *ast.ImportSpec) string { return t.Path.Value })).
			Map(func(s string) string {
				return strings.Trim(s, "\"")
			}).
			Filter(func(s string) bool {
				return len(stream.OfSlice(strings.Split(s, "/")).
					Filter(func(ss string) bool {
						return ss == alias
					}).ToSlice()) > 0
			}).One()
		if imp != "" {
			return imp, true
		}
		return "", false
	}
}

func GetImports(e ast.Expr, lookup ImportLookup) DistinctImports {
	out := NewDistinctImports()
	ast.Inspect(e, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ArrayType:
			i := GetImports(getSelectorExpr(n), lookup)
			out.Merge(i)
		case *ast.SelectorExpr:
			if n == nil || n.X == nil {
				return false
			}
			switch i := n.X.(type) {
			case *ast.Ident:
				alias := i.String()
				pkg, ok := lookup(alias)
				if !ok {
					return true
				}
				out.Append(Import{
					Alias:   alias,
					Package: pkg,
				})
			}
		}
		return true
	})
	return out
}
