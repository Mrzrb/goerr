package aop_core

import "encoding/json"

type Jointcut struct {
	TargetName string
	TargetType string
	MethodName string
	// args
	Args []Args
	// warp process
	Fn       func() error
	MeteInfo string
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

func (j *Jointcut) GetMetaInfo() *JointMete {
	var m JointMete

	err := json.Unmarshal([]byte(j.MeteInfo), &m)
	if err != nil {
		return nil
	}

	return &m
}

type Mete struct {
	Mete map[string]map[string]string `json:"mete"`
}

type JointMete struct {
	StructMeta    Mete            `json:"structMeta"`
	ProcedureMeta map[string]Mete `json:"procedureMeta"`
}
