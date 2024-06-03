package autowire

import (
	"errors"
	"fmt"
	"go/ast"
	"os"
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
	f := utils.GetFullPackage(".").Module.Dir + "/" + a.BaseOutputer.File
	fmt.Println(a.BaseOutputer.File)
	fslice := strings.Split(f, "/")
	f = strings.Join(fslice[:len(fslice)-1], "/")
	// 创建文件路径中可能不存在的目录
	os.MkdirAll(f, 0755)
	return a.BaseOutputer.File
}

func fieldAutowire(field *ast.Field) bool {
	if field.Doc == nil {
		return false
	}
	for _, comment := range field.Doc.List {
		if strings.Contains(comment.Text, "Autowire") {
			return true
		}
	}

	return false
}

// Imports implements core.Outputer.
func (p *Assembler) Imports() []string {
	for _, v := range p.Components {
		if v.Meta().PackageName() != p.Package() {
			p.BaseOutputer.Imports = append(p.BaseOutputer.Imports, utils.GetPkgWithFullPath(v.Meta()))
		}
		v.WalkField(func(f *ast.Field) {
			if fieldAutowire(f) {
				p.BaseOutputer.Imports = append(p.BaseOutputer.Imports, v.Node.Import(f.Type)...)
			}
		})
	}
	for _, v := range p.Factory {
		if v.Meta().PackageName() != p.Package() {
			p.BaseOutputer.Imports = append(p.BaseOutputer.Imports, utils.GetPkgWithFullPath(v.Meta()))
		}
		v.WalkField(func(f *ast.Field) {
			if fieldAutowire(f) {
				p.BaseOutputer.Imports = append(p.BaseOutputer.Imports, v.Node.Import(f.Type)...)
			}
		})
	}
	return append(p.BaseOutputer.Imports, "github.com/samber/do/v2", "github.com/Mrzrb/goerr/di", "fmt")
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
		vType := utils.OrGet(v.Meta().PackageName() == "main", v.Type, v.Meta().PackageName()+"."+v.Type)
		tp := utils.OrGet(declIsPointer, "*"+vType, vType)
		p.b.Write(GetRegisterVariable(decl, v.Id(), tp, variableName))
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
			val := initial.Name
			if f.IsPointer && !initial.IsPoint {
				val = "&" + initial.Name
			}
			if !f.IsPointer && initial.IsPoint {
				val = "*" + initial.Name
			}
			return fmt.Sprintf("%s: %s,", t.Name, val)
		})
		decl := GetInitVariable(variableName, utils.OrGet(v.Meta().PackageName() == "main", v.Type, v.Meta().PackageName()+"."+v.Type), params)
		p.b.Write(GetRegisterVariable(string(decl), v.Id(), utils.OrGet(v.Meta().PackageName() == "main", v.Type, v.Meta().PackageName()+"."+v.Type), variableName))
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
