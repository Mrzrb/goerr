package autowire

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
)

type PreCollector struct {
	Components []string
	Factories  []string
}

type Assembler struct {
	Scope      string
	Pre        PreCollector
	Components []*Struct
	Factory    []*core.Func
	Fields     map[string]*core.Field
	core.BaseOutputer
	b strings.Builder

	// if one component initial, register here with Id(), and it variableName
	inited map[string]initItem
}

type initItem struct {
	IsPoint bool
	Name    string
}

func NewInitItem(isPoint bool, name string) initItem {
	return initItem{
		IsPoint: isPoint,
		Name:    name,
	}
}

// File implements core.Outputer.
func (a *Assembler) File() string {
	return "autowire.gen.go"
}

// Imports implements core.Outputer.
func (a *Assembler) Imports() []string {
	return append(a.BaseOutputer.Imports, "github.com/samber/do/v2", "github.com/Mrzrb/goerr/di")
}

// Output implements core.Outputer.
func (a *Assembler) Output() []byte {
	return []byte(warp(a.Inject()))
}

// Package implements core.Outputer.
func (a *Assembler) Package() string {
	return a.BaseOutputer.Package
}

// Valid implements core.Outputer.
func (a *Assembler) Valid() error {
	return nil
}

type Struct struct {
	core.Struct
	Depend DependencyTree
}

type DependencyTree struct {
	Id                string
	Name              string
	ChildDependencies []*DependencyTree
	Resolved          bool
}

func (a *Assembler) GetGenPackage() string {
	ps := strings.Split(a.BaseOutputer.File, "/")
	if len(ps) <= 1 {
		panic("gen filename invalid")
	}
	return ps[len(ps)-1]
}

var _ core.Outputer = (*Assembler)(nil)

func (p *Process) GetUnit() *Assembler {
	ret := &Assembler{
		Scope:        "",
		Pre:          PreCollector{},
		Components:   []*Struct{},
		Factory:      []*core.Func{},
		BaseOutputer: core.BaseOutputer{},
	}

	return ret
}

func (p *Assembler) parseComponent(v *Struct) error {
	name := v.Name
	fmt.Println(name)
	if _, ok := p.inited[v.Id()]; ok {
		return nil
	}
	canInitial := true
	for _, f := range v.Depend.ChildDependencies {
		if !f.Resolved {
			canInitial = false
			fieldStruct := p.GetComponentById(f.Id)
			err := p.parseComponent(fieldStruct)
			if err == nil {
				canInitial = true
			}
		}
	}
	if !canInitial {
		return errors.New("cannot initial")
	}
	c := p.GetFactoryById(v.Id())
	// in has factory, get instance from factory
	declIsPointer := false
	if c != nil {
		if len(c.Retern) > 1 {
			panic("factory must be one")
		}
		declIsPointer = c.Retern[0].IsPointer
		variableRet := strings.ReplaceAll(c.Id(), ".", "_")
		variableRet = strings.ReplaceAll(variableRet, "/", "_")
		variableRet = strings.ToLower(variableRet)
		decl, variableName := GetTplFactory(c, p.BaseOutputer.Package, "", []string{variableRet})
		p.b.Write(GetRegisterVariable(decl, v.Id(), utils.OrGet(declIsPointer, "*"+v.Type, v.Type), variableName))
		p.inited[v.Id()] = NewInitItem(declIsPointer, variableName)
	} else {
		variableName := strings.ReplaceAll(v.Id(), ".", "_")
		variableName = strings.ReplaceAll(variableName, "/", "_")
		if p.BaseOutputer.Package == v.Meta().PackageName() {
			variableName = v.Name
		}
		variableName = strings.ToLower(variableName)
		params := utils.Map(v.Depend.ChildDependencies, func(t *DependencyTree) string {
			initial, ok := p.inited[t.Id]
			if !ok {
				panic(fmt.Sprintf("initial fail, component not init, %s", t.Id))
			}
			f, ok := p.Fields[t.Id]
			if !ok {
				panic("cannot find param")
			}
			val := ""
			if f.IsPointer && initial.IsPoint {
				val = initial.Name
			}
			if f.IsPointer && !initial.IsPoint {
				val = "&" + initial.Name
			}
			if !f.IsPointer && initial.IsPoint {
				val = "*" + initial.Name
			}
			return fmt.Sprintf("%s: %s,", t.Name, val)
		})
		decl := GetInitVariable(variableName, v.Type, params)
		p.b.Write(GetRegisterVariable(string(decl), v.Id(), v.Type, variableName))
		p.inited[v.Id()] = NewInitItem(false, variableName)
	}
	return nil
}

func (p *Assembler) Inject() []byte {
	for _, v := range p.Components {
		p.parseComponent(v)
	}
	return []byte(p.b.String())
}

func (p *Assembler) GetFactoryById(id string) *core.Func {
	for _, c := range p.Factory {
		for _, r := range c.FuncIdent.Retern {
			if r.Id(c.Node) == id {
				return c
			}
		}
	}
	return nil
}

func (p *Assembler) GetComponentById(id string) *Struct {
	for _, c := range p.Components {
		if c.Id() == id {
			return c
		}
	}
	return nil
}