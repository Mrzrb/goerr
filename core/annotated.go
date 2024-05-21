package core

import annotation "github.com/YReshetko/go-annotation/pkg"

type Annotated interface {
	Annotate() []annotation.Annotation
}

type AnnotationsMix struct {
	Annotation []annotation.Annotation
}

type Ident struct {
	AnnotationsMix
	Name string
	Type string
}

type FuncIdent struct {
	AnnotationsMix
	Name   string
	Param  []Ident
	Retern []Ident
}

func Cast[T Annotated](n Annotated) (T, bool) {
	v, ok := n.(T)
	return v, ok
}
