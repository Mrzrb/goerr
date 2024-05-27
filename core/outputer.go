package core

import (
	"fmt"
	"strings"

	"github.com/Mrzrb/goerr/utils"
)

// @Mock
type Outputer interface {
	outputer
	File() string
}

type outputer interface {
	Valid() error
	Output() []byte
	Imports() []string
	Package() string
}

type FileExporter interface {
	Append(Outputer)
	Export() map[string][]byte
}

// @Getter
type BaseOutputer struct {
	File    string
	Imports []string
	Package string
}

func (b *BaseOutputer) PushImport(pkgName string, dirName string) {
	pkg := pkgName
	if pkgName != b.Package {
		pkg = utils.GetFullPackage(dirName)
	}

	b.Imports = append(b.Imports, pkg)
}

type BaseFuncOutputer struct {
	IsMethod  bool
	Receiver  string
	FuncName  string
	Params    []Ident
	Param     string
	Returns   []Ident
	Return    string
	CallParam string
	HasReturn bool
}

func (b *BaseFuncOutputer) AssembleParamString() string {
	return strings.Join(utils.Map(b.Params, func(t Ident) string {
		return fmt.Sprintf("%s %s", t.Name, t.Type)
	}), ",")
}

func (b *BaseFuncOutputer) AssembleCallParamString() string {
	return strings.Join(utils.Map(b.Params, func(t Ident) string {
		return fmt.Sprintf("%s", t.Name)
	}), ",")
}

func (b *BaseFuncOutputer) AssembleReturnString() string {
	idx := 0
	return "(" + strings.Join(utils.Map(b.Returns, func(t Ident) string {
		idx++
		return fmt.Sprintf("ret%d %s", idx, t.Type)
	}), ",") + ")"
}

func (b *BaseFuncOutputer) AssembleReturnResultAppendString() []string {
	idx := 0
	return utils.Map(b.Returns, func(t Ident) string {
		idx++
		return fmt.Sprintf(`returnResult.Args = append(returnResult.Args, &aop.Args{Name:"%s",Type:"%s",Value:ret%d})`, t.Name, t.Type, idx)
		// return fmt.Sprintf("ret%d %s", idx, t.Type)
	})
}

func (b *BaseFuncOutputer) AssembleReturnDecl() string {
	idx := 0
	return strings.Join(utils.Map(b.Returns, func(t Ident) string {
		idx++
		return fmt.Sprintf("ret%d", idx)
	}), ",")
}

func (b *BaseFuncOutputer) AssembleResultSetString() []string {
	idx := 0
	return utils.Map(b.Returns, func(t Ident) string {
		s := fmt.Sprintf(`runContext.ReturnResult.Args[%d].Value = ret%d`, idx, idx+1)
		idx++
		return s
	})
}

func (b *BaseFuncOutputer) AssembleErrorCheckers() []string {
	idx := 0
	return utils.Map(b.Returns, func(t Ident) string {
		idx++
		return utils.OrGet(t.Type == "error", fmt.Sprintf(`if "%s" == "error" && ret%d != nil {
            return ret%d
        }`, t.Type, idx, idx), "")
		// return fmt.Sprintf(`if "%s" == "error" {
		//           %s
		//       }`, t.Type, utils.OrGet(t.Type == "error", fmt.Sprintf("return ret%d", idx), ""))
	})
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

var Assembler = NewAssembler()

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
			"Import": utils.Uniq(o.Imports(), func(t string) string { return t }),
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
				e.cache[f] = 1
			}
			if err := v.Valid(); err != nil {
				panic(fmt.Sprintf("param not valid %s", err))
			}
			ret[f] = append(ret[f], v.Output()...)
		}
	}

	return ret
}

var _ FileExporter = (*exporter)(nil)
