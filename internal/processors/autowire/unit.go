package autowire

import (
	"strings"

	"github.com/Mrzrb/goerr/core"
)

type PreCollector struct {
	Components []string
	Factories  []string
}

type Assembler struct {
	File       string
	Scope      string
	Pre        PreCollector
	Components []*Struct
	Factory    []*core.Func
	walked     []string
}

type Struct struct {
	core.Struct
	Depend DependencyTree
}

type DependencyTree struct {
	Id                string
	ChildDependencies []*DependencyTree
}

func (a *Assembler) GetGenPackage() string {
	ps := strings.Split(a.File, "/")
	if len(ps) <= 1 {
		panic("gen filename invalid")
	}
	return ps[len(ps)-1]
}

type Demo struct{}

func init() {
}
