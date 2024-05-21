package core

import (
	"bytes"
	"fmt"
	"html/template"

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

const (
	pkgName    = "pkgName"
	importName = "import"
)

var temp *template.Template

func init() {
	temp = utils.Must(template.New(pkgName).Parse(tpl))
	temp = utils.Must(temp.New(importName).Parse(imprtTpl))
}

func ExecuteTemplate(tpl *template.Template, data any) ([]byte, error) {
	b := bytes.NewBufferString("")
	err := tpl.Execute(b, data)
	if err != nil {
		return nil, fmt.Errorf("unable to process template %s: %w", tpl.Name(), err)
	}
	return b.Bytes(), nil
}
