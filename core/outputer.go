package core

import (
	"fmt"

	"github.com/Mrzrb/goerr/utils"
)

// @Mock
type Outputer interface {
	outputer
	File() string
}

type outputer interface {
	Valid() bool
	Output() []byte
	Imports() []string
	Package() string
}

type FileExporter interface {
	Append(Outputer)
	Export() map[string][]byte
}

// @Constructor
type exporter struct {
	Files     []Outputer // @Init
	Assembler *assembler
	cache     map[string]any // @Init
}

// Append implements FileExporter.
func (e *exporter) Append(f Outputer) {
	e.Files = append(e.Files, f)
}

// @Constructor
type assembler struct{}

func (a *assembler) Assmble(o outputer) []byte {
	ret := []byte{}
	pkg := temp.Lookup(pkgName)
	if pkg == nil {
		panic(pkg)
	}
	ret = append(ret, utils.Must(ExecuteTemplate(pkg, map[string]any{
		"PackageName": o.Package(),
	}))...)
	if len(o.Imports()) > 0 {
		imp := temp.Lookup(importName)
		if imp == nil {
			panic(pkg)
		}
		ret = append(ret, utils.Must(ExecuteTemplate(imp, map[string]any{
			"Import": o.Imports(),
		}))...)
	}
	return ret
}

// export implements FileExporter.
func (e *exporter) Export() map[string][]byte {
	ret := map[string][]byte{}
	groupedExporter := utils.GroupBy(e.Files, func(t Outputer) string {
		return t.File()
	})

	for f, ge := range groupedExporter {
		for _, v := range ge {
			if _, ok := e.cache[f]; !ok {
				ret[f] = append(ret[f], e.Assembler.Assmble(v)...)
			}
			if !v.Valid() {
				panic(fmt.Sprintf("param not valid %+v", v))
			}
			ret[f] = append(ret[f], v.Output()...)
		}
	}

	return ret
}

var _ FileExporter = (*exporter)(nil)
