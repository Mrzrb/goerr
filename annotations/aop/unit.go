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
	Effects    []Effect
	AspectType string
	Method     []core.BaseFuncOutputer
}

type Effect struct {
	Aspect     core.Struct
	Affect     core.Method
	AspectType string
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
		"Name": u.Name,
		"Type": u.Type,
		"AspectTypeDecl": utils.MapIdx(u.Effects, func(t Effect, index int) string {
			return "aspect" + fmt.Sprintf("%d", index) + " *" + utils.OrGet(u.GetPackage() == t.Aspect.Meta().PackageName(), t.Aspect.Type, t.Aspect.Meta().PackageName()+"."+t.Aspect.Type)
		}),
		"AspectTypeDeclInit": utils.MapIdx(u.Effects, func(t Effect, index int) string {
			return "aspect" + fmt.Sprintf("%d", index) + ": &" + utils.OrGet(u.GetPackage() == t.Aspect.Meta().PackageName(), t.Aspect.Type, t.Aspect.Meta().PackageName()+"."+t.Aspect.Type) + "{}"
		}),
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
			"Effects":            u.Effects,
			"CallJoints": utils.MapIdx(u.Effects, func(t Effect, index int) string {
				return fmt.Sprintf("r.aspect%d.%s(j, m)", index, t.Affect.Name)
			}),
		}
		m["Return"] = v.AssembleReturnString()
		m["CallParams"] = v.AssembleCallParamString()
		m["ReturnVal"] = v.AssembleReturnDecl()
		m["ErrorCheckers"] = v.AssembleErrorCheckers()
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
func (u *Unit) Valid() error {
	for _, v := range u.Effects {
		if len(v.Affect.Retern) != 1 {
			return fmt.Errorf("affect must return error, Aspect %s", v.Aspect.Name)
		}
		for _, vv := range v.Affect.Retern {
			if vv.Type != "error" {
				return fmt.Errorf("affect must return error, actual return %s, Aspect %s", vv.Type, v.Aspect.Name)
			}
		}
	}
	return nil
}

var _ core.Outputer = (*Unit)(nil)
