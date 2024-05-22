package logger

import (
	"text/template"

	"github.com/Mrzrb/goerr/utils"
)

var tpl = `
func {{if .IsMethod}}(r {{.Receiver}}){{end}} {{.FuncName}}Logger({{.Param}}){{.Return}} {
    now := time.Now().Unix()
    zlog.Infof(ctx, "{{.FuncName}} enter_time %d", now)
    defer func() {
        now1 := time.Now().Unix()
        zlog.Infof(ctx, "{{.FuncName}} exit_time %d, duration %d", now1, now1-now)
    }()
    {{if .HasReturn}}
    return {{if .IsMethod}}r.{{end}}{{.FuncName}}({{.CallParam}})
    {{else}}
    {{if .IsMethod}}r.{{end}}{{.FuncName}}({{.CallParam}})
    return
    {{end}}
}
`

const tplName = "tpl"

var temp *template.Template

func init() {
	temp = utils.Must(template.New(tplName).Parse(tpl))
}
