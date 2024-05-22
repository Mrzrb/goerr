package aop

var tpl = `
type {{.Name}}Proxy struct {
    inner *{{.Type}}
    aspect *{{.AspectType}}
}

func New{{.Name}}Proxy(inner *{{Type}}) {{.Name}}Proxy {
    return &{{.Name}}Proxy {
        inner: inener,
        aspect: &{{.AspectType}}{},
    }
}

func (r *{{Type}}) {{FuncName}}({{.Param}}) {{.Return}} {
    {{range .Returns}}
    var {{.Name}} {{.Type}}{{end}}
	joint := aop.Jointcut{
		TargetName: "{{.Name}}",
		TargetType: "{{.Type}}",
		Args:       []aop.Args{},
		Fn: func() {
			ret1 = r.inner.{{.FuncName}}({{.CallParams}})
		},
	}
    {{range .Params}}
    joint.Args = append(joint.Args, aop.Args{ {{.Name}}, {{.Type}} }){{end}}

}
`
