package aop

import "github.com/Mrzrb/goerr/core"

type Aop struct {
	Type   string `annotation:"name=type,default=aspect,oneOf=aspect;point;pointcut;before;after;around;catchPanic"`
	Target string `annotation:"name=target"`
}

type MixinType string

const (
	Before     MixinType = "before"
	Around     MixinType = "around"
	After      MixinType = "after"
	CatchPanic MixinType = "catchPanic"
)

type Jointcut struct {
	TargetName string
	TargetType string
	MethodName string
	// args
	Args []Args
	// warp process
	Fn func() error
}

type Args struct {
	Name  string
	Type  string
	Value any
}

func IdentToArgs(i core.Ident) Args {
	return Args{
		Name: i.Name,
		Type: i.Type,
	}
}
