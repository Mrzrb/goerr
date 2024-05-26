package aop

import (
	"text/template"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

var tpl = `
type {{.Name}}Proxy struct {
    inner *{{.Type}}
    {{range .AspectTypeDecl}}
    {{.}}{{end}}
}

func New{{.Name}}Proxy(inner *{{.Type}}) *{{.Name}}Proxy {
    return &{{.Name}}Proxy {
        inner: inner,
        {{range .AspectTypeDeclInit}}
        {{.}},{{end}}
    }
}
`

var tplInterface = `
type {{.Type}}Interface interface {
{{range .Methods}}
    {{.FuncName}} ({{.Param}}) {{.Return}}
{{end}}
}
`

var tplMethod = `
{{range .Methods}}
func (r *{{.Name}}Proxy) {{.FuncName}}({{.Param}}) {{.Return}} {
	joint := aop.Jointcut{
		TargetName: "{{.Name}}",
		TargetType: "{{.Type}}",
		Args:       []aop.Args{},
		Fn: func() {
            {{.ReturnVal}} = r.inner.{{.FuncName}}({{.CallParams}})
		},
	}
    {{range .Params}}
    joint.Args = append(joint.Args, aop.Args{ Name : "{{.Name}}", Type: "{{.Type}}", Value: {{.Name}} }){{end}}

    fn := aop.GenerateChain(joint,
        {{range .CallJoints}}
        func(j aop.Jointcut) {
            {{.}}
        },
        {{end}}
    )
    fn()
    return {{.ReturnVal}}
}{{end}}
`

var (
	tplName          = "proxy"
	tplMethodName    = "method"
	tplInterfaceName = "interface"
)

var temp *template.Template

func init() {
	temp = utils.Must(template.New(tplName).Funcs(core.CustomeTemplateFuncs()).Parse(tpl))
	temp = utils.Must(temp.New(tplMethodName).Funcs(core.CustomeTemplateFuncs()).Parse(tplMethod))
	temp = utils.Must(temp.New(tplInterfaceName).Funcs(core.CustomeTemplateFuncs()).Parse(tplInterface))
}
