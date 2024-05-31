package aop

import (
	"github.com/Mrzrb/goerr/annotations/aop_core"
	"github.com/Mrzrb/goerr/core"
)

type Aop struct {
	Type   string `annotation:"name=type,default=aspect,oneOf=aspect;point;pointcut;before;after;around;catchPanic"`
	Target string `annotation:"name=target"`
	Val    string `annotation:"name=val"`
}

type MixinType string

const (
	Before     MixinType = "before"
	Around     MixinType = "around"
	After      MixinType = "after"
	CatchPanic MixinType = "catchPanic"
)

func IdentToArgs(i core.Ident) aop_core.Args {
	return aop_core.Args{
		Name: i.Name,
		Type: i.Type,
	}
}
