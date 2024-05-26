package aop

import (
	"fmt"
	"strings"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

type Unit struct {
	core.BaseOutputer
	core.Struct
	Aspect     core.Struct
	Affect     core.Method
	AspectType string
	Method     []core.BaseFuncOutputer
}

// File implements core.Outputer.
func (u *Unit) File() string {
	return u.GetFile()
}

// Imports implements core.Outputer.
func (u *Unit) Imports() []string {
	return append(u.GetImports(), "github.com/Mrzrb/goerr/annotations/aop")
}

// Output implements core.Outputer.
func (u *Unit) Output() []byte {
	ret := []byte{}
	t := utils.MustPointer(temp.Lookup(tplName))
	ret = append(ret, utils.Must(core.ExecuteTemplate(t, map[string]any{
		"Name":       u.Name,
		"Type":       u.Type,
		"AspectType": utils.OrGet(u.GetPackage() == u.Aspect.Meta().PackageName(), u.Aspect.Type, u.Aspect.Meta().PackageName()+"."+u.Aspect.Type),
	}))...)

	tM := utils.MustPointer(temp.Lookup(tplMethodName))
	tI := utils.MustPointer(temp.Lookup(tplInterfaceName))

	param := map[string]any{
		"Name":    u.Name,
		"Type":    u.Type,
		"Methods": []map[string]any{},
	}
	for _, v := range u.Method {
		m := map[string]any{
			"Name":     u.Name,
			"Type":     u.Type,
			"FuncName": v.FuncName,
			"Param":    v.AssembleParamString(),
			"ReturnDecl": strings.Join(utils.Map(v.Returns, func(t core.Ident) string {
				return fmt.Sprintf("var %s %s", t.Name, t.Type)
			}), "\n"),
			"Params":             v.Params,
			"AffectedMethodName": u.AspectType,
			"ReturnDeclNames":    v.Return,
		}
		m["Return"] = v.AssembleReturnString()
		m["CallParams"] = v.AssembleCallParamString()
		m["ReturnVal"] = v.AssembleReturnDecl()
		param["Methods"] = append(param["Methods"].([]map[string]any), m)
	}

	ret = append(ret, utils.Must(core.ExecuteTemplate(tI, param))...)
	ret = append(ret, utils.Must(core.ExecuteTemplate(tM, param))...)

	return ret
}

// Package implements core.Outputer.
func (u *Unit) Package() string {
	return u.GetPackage()
}

// Valid implements core.Outputer.
func (u *Unit) Valid() bool {
	return true
}

var _ core.Outputer = (*Unit)(nil)
