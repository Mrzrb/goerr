package autowire

import (
	"fmt"
	"text/template"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

func warp(body []byte) string {
	fn := `
func init() {
    %s
}
    `
	return fmt.Sprintf(fn, string(body))
}

var tplRegisterVariable = `
    {{.Declare}}
    do.ProvideNamedValue(di.GlobalInjector, "{{.InjectName}}", {{.InjectVal}})
    do.Provide(di.GlobalInjector, func(i do.Injector) ({{.InjectType}}, error) {
        return {{.InjectVal}}, nil
    })
`
var tplRegisterVariableName = "register"

var tplInitVariable = `
    {{.InjectVal}} := {{.InjectType}}{
        {{range .Params}}{{.}}
        {{end}}
    }
`
var tplInitVariableName = "init"

var tpl *template.Template

func init() {
	tpl = utils.Must(template.New(tplRegisterVariableName).Parse(tplRegisterVariable))
	tpl = utils.Must(tpl.New(tplInitVariableName).Parse(tplInitVariable))
}

func GetInitVariable(injectVal string, injectType string, p []string) []byte {
	params := map[string]any{
		"InjectVal":  injectVal,
		"InjectType": injectType,
		"Params":     p,
	}
	t := utils.MustPointer(tpl.Lookup(tplInitVariableName))
	return utils.Must(core.ExecuteTemplate(t, params))
}

func GetRegisterVariable(declare string, variableName string, injectType string, injectVal string) []byte {
	params := map[string]any{
		"Declare":    declare,
		"InjectName": variableName,
		"InjectType": injectType,
		"InjectVal":  injectVal,
	}
	t := utils.MustPointer(tpl.Lookup(tplRegisterVariableName))
	bt, err := core.ExecuteTemplate(t, params)
	if err != nil {
		panic(err)
	}
	return bt
}

func GetTplFactory(c *core.Func, pkg string, receiver string, returns []string, params ...string) (decl string, variableName string) {
	decl = c.Call(pkg, receiver, returns, params...)
	if len(returns) > 0 {
		variableName = returns[0]
	}

	return
}
