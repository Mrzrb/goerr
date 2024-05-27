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

type RunContext struct {
	MuteableArgs
	ReturnResult
}

type MuteableArgs struct {
	Args []*Args
}

type ReturnResult struct {
	Args []*Args
}

func (j *Jointcut) Copy() Jointcut {
	return Jointcut{
		TargetName: j.TargetName,
		TargetType: j.TargetType,
		MethodName: j.MethodName,
		Args:       j.Args,
		Fn: func() error {
			return nil
		},
	}
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
