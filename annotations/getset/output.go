package getset

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
	"text/template"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

var tpl = `
package {{.PackageName}}
`

var imprtTpl = `
import (
    {{range .Import}} "{{ . }}"  {{end}}
)
`

var getterTpl = `
func (c {{ .Receiver }}) Get{{.Name}}() {{.Type}} {
    return c.{{.RawName}}
}
`

var setterTpl = `
func (c {{ .Receiver }}) Set{{.Name }}(val {{ .Type }}) {{.Receiver}} {
    c.{{.RawName}} = val
    return c
}
`

const (
	pkgName    = "pkgName"
	importName = "import"
	setterName = "setter"
	getterName = "getter"
)

var oIns out = out{}

type out struct {
	temp *template.Template
	g    *GsProcessor
	// out  map[string]outStruct
}

func init() {
	temp := utils.Must(template.New(pkgName).Parse(tpl))
	temp = utils.Must(temp.New(getterName).Parse(getterTpl))
	temp = utils.Must(temp.New(setterName).Parse(setterTpl))
	temp = utils.Must(temp.New(importName).Parse(imprtTpl))
	oIns.temp = temp
}

func (o *out) Output(g *GsProcessor) map[string][]byte {
	ret := map[string][]byte{}
	groupParsed := utils.GroupBy(g.Parsed, func(t core.Annotated) string {
		return t.Nodes().Meta().FileName()
	})
	for _, v := range groupParsed {
		parsed := utils.Map(v, func(t core.Annotated) *core.Struct {
			s := utils.MustBool(core.Cast[*core.Struct](t))
			return s
		})
		dstFile, o := o.output(parsed)
		ret[dstFile] = o
	}
	return ret
}

func (o *out) output(s []*core.Struct) (string, []byte) {
	ret := []byte{}
	if len(s) == 0 {
		return "", nil
	}
	ret = append(ret, o.parsePackage(s[0].Meta().PackageName())...)

	imprts := []string{}
	for _, i := range s {
		for _, f := range i.Field {
			im := utils.GetImports(f.Raw.(*ast.Field).Type, i.Lookup().FindImportByAlias)
			if len(im.ToSlice()) > 0 {
				imprts = append(imprts, im.ToSlice()[0].Package)
			}
		}
	}
	if len(imprts) > 0 {
		ret = append(ret, o.parseImport(imprts)...)
	}

	for _, parsed := range s {
		s := utils.MustBool(core.Cast[*core.Struct](parsed))
		if core.ContainsAnnotate[GetterSetter](s) {
			ret = append(ret, o.parseGetter(s)...)
			ret = append(ret, o.parseSetter(s)...)
			continue
		}
		if core.ContainsAnnotate[Getter](s) {
			ret = append(ret, o.parseGetter(s)...)
		}
		if core.ContainsAnnotate[Setter](s) {
			ret = append(ret, o.parseSetter(s)...)
		}
	}

	f := s[0].DstFileName()
	return f, ret
}

func (o *out) parsePackage(pkg string) []byte {
	pkgTpl := o.temp.Lookup(pkgName)
	if pkgTpl != nil {
		return utils.Must(ExecuteTemplate(pkgTpl, map[string]any{
			"PackageName": pkg,
		}))
	}
	return nil
}

func (o *out) parseImport(imprts []string) []byte {
	impTpl := o.temp.Lookup(importName)
	if impTpl != nil {
		return utils.Must(ExecuteTemplate(impTpl, map[string]any{
			"Import": imprts,
		}))
	}
	return nil
}

func (o *out) parseGetter(data *core.Struct) []byte {
	ret := []byte{}
	getterTpl := o.temp.Lookup(getterName)
	for _, v := range data.Field {
		if v.Name == "" {
			continue
		}
		param := map[string]any{
			"Receiver": "*" + data.Type,
			"Name":     strings.Title(v.Name),
			"RawName":  v.Name,
			"Type":     v.Type,
		}
		ret = append(ret, utils.Must(ExecuteTemplate(getterTpl, param))...)
	}

	return ret
}

func (o *out) parseSetter(data *core.Struct) []byte {
	ret := []byte{}
	setterTpl := o.temp.Lookup(setterName)
	for _, v := range data.Field {
		if v.Name == "" {
			continue
		}
		param := map[string]any{
			"Receiver": "*" + data.Type,
			"Name":     strings.Title(v.Name),
			"RawName":  v.Name,
			"Type":     v.Type,
		}
		ret = append(ret, utils.Must(ExecuteTemplate(setterTpl, param))...)
	}

	return ret
}

func ExecuteTemplate(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	return b.Bytes(), nil
}

func (o *out) Out() map[string][]byte {
	return o.Output(o.g)
}
