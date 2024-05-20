package utils

import (
	"go/ast"
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

func GetImports(e ast.Expr, lookup ImportLookup) DistinctImports {
	out := NewDistinctImports()
	ast.Inspect(e, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.SelectorExpr:
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
