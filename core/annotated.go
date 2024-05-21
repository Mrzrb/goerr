package core

import annotation "github.com/YReshetko/go-annotation/pkg"

type Annotated interface {
	Annotate() []annotation.Annotation
}

type Ident struct {
	Name       string
	Type       string
	Annotation []annotation.Annotation
}

func Cast[T Annotated](n Annotated) (T, bool) {
	v, ok := n.(T)
	return v, ok
}
