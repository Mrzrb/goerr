package aop_core

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
