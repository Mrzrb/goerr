package getset

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Mrzrb/goerr/utils"
)

var tpl = `
package {{.PackageName}}
`

var imprtTpl = `
    import (
        {{range .Import}} {{ . }}  {{end}}
    )
`

var getterTpl = `
func (c *{{ .TargetName }}) Get{{.Name}}() {{.Type}} {
    return c.{{.Name}}
}
`

var setterTpl = `
func (c *{{ .TargetName }}) Set{{.Name }}(val {{ .Type }}) *{{.TargetName}} {
    c.{{.Name}} = val
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
	out  map[string]outStruct
}

type outStruct struct {
	Import       []string
	PackageName  string
	fullFileName string
	Fields       []gsTarget
}

func ExecuteTemplate(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	return b.Bytes(), nil
}

func (os *outStruct) Out(formatter *template.Template) []byte {
	ret := []byte{}
	pkgTpl := formatter.Lookup(pkgName)
	if pkgTpl != nil {
		ret = append(ret, utils.Must(ExecuteTemplate(pkgTpl, os))...)
	} else {
		return nil
	}

	if len(os.Import) > 0 {
		impTpl := formatter.Lookup(importName)
		if impTpl != nil {
			ret = append(ret, utils.Must(ExecuteTemplate(impTpl, os))...)
		}
	}

	if len(os.Fields) > 0 {
		setterTpl := formatter.Lookup(setterName)
		getterTpl := formatter.Lookup(getterName)
		for _, f := range os.Fields {
			for _, ff := range f.field {
				if ff.IsGetter {
					ret = append(ret, utils.Must(ExecuteTemplate(getterTpl, ff))...)
				}
				if ff.IsSetter {
					ret = append(ret, utils.Must(ExecuteTemplate(setterTpl, ff))...)
				}
			}
		}
	}
	return ret
}

func init() {
	temp := utils.Must(template.New(pkgName).Parse(tpl))
	temp = utils.Must(temp.New(getterName).Parse(getterTpl))
	temp = utils.Must(temp.New(setterName).Parse(setterTpl))
	temp = utils.Must(temp.New(importName).Parse(imprtTpl))
	oIns.temp = temp
}

func (o *out) Out() map[string][]byte {
	ret := map[string][]byte{}
	for file, outStruct := range o.out {
		fnames := strings.Split(file, ".")
		utils.InsertInPos(&fnames, "gen", int64(len(fnames)-1))
		ret[strings.Join(fnames, ".")] = outStruct.Out(o.temp)
	}
	return ret
}

func (o *out) Process(gs *GsProcessor) {
	for _, v := range gs.Targets {
		fname := v.FullFileName()
		outS, ok := o.out[fname]
		if !ok {
			os := outStruct{
				Import:       []string{},
				PackageName:  v.packageName,
				fullFileName: fname,
				Fields:       []gsTarget{},
			}
			o.out = map[string]outStruct{}
			o.out[fname] = os
			outS = os
		}
		if len(v.Imps) > 0 {
			imports := []string{}
			for distinctImport := range v.Imps {
				imports = append(imports, fmt.Sprintf(`"%s"`, distinctImport.Package))
				// if distinctImport.Alias != "" {
				// 	imports = append(imports, fmt.Sprintf(`%s "%s"`, distinctImport.Alias, distinctImport.Package))
				// } else {
				// }
			}
			outS.Import = append(outS.Import, imports...)
		}
		outS.Fields = append(outS.Fields, v)
		o.out[fname] = outS
	}
}
