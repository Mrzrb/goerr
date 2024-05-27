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
        MethodName: "{{.FuncName}}",
		Args:       []aop.Args{},
	}
    {{range .Params}}
    joint.Args = append(joint.Args, aop.Args{ Name : "{{.Name}}", Type: "{{.Type}}", Value: {{.Name}} }){{end}}

    runContext := aop.RunContext{}
    returnResult := aop.ReturnResult{}
    {{range .ResultAppend}}
    {{.}}{{end}}

    mutableArgs := aop.MuteableArgs{}
    {{range $idx, $e := .Params}}
    mutableArgs.Args = append(mutableArgs.Args, &joint.Args[{{$idx}}]){{end}}
    runContext.MuteableArgs = mutableArgs
    runContext.ReturnResult = returnResult

    
    joint.Fn = func() error {
            {{.ReturnValSet}} = r.inner.{{.FuncName}}({{range $idx, $e := .Params}}mutableArgs.Args[{{$idx}}].Value.({{.Type}}),{{end}})
            {{range .ErrorCheckers}}
            {{.}}{{end}}
            {{range .ResultSet}}
            {{.}}{{end}}
            return nil
    }

    aop.GenerateChain(&joint,&runContext,
        {{range .CallJoints}}
        func(j aop.Jointcut, m *aop.RunContext) error {
            return {{.}}
        },
        {{end}}
    )
    joint.Fn()
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
