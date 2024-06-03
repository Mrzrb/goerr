package core

import (
	"fmt"
	"go/ast"
	"reflect"

	annotation "github.com/YReshetko/go-annotation/pkg"
)

type Field struct {
	Ident
	Package         string
	Alias           string
	FullPackagePath string
	Parent          annotation.Node
}

// Id implements Identity.
func (f *Field) Id() string {
	parentNode := f.Parent
	if parentNode == nil {
		panic("field must have parent struct")
	}
	pNode := Parse(parentNode)

	if _, ok := Cast[*Struct](pNode); ok {
		_, impath, _ := findFieldPackage(f.Name, parentNode)
		return fmt.Sprintf("%s.%s", impath, f.Type)
	}

	return ""
}

// Annotate implements Annotated.
func (f *Field) Annotate() []annotation.Annotation {
	return f.Annotate()
}

// DstFileName implements Annotated.
func (f *Field) DstFileName(...string) string {
	return ""
}

// Nodes implements Annotated.
func (f *Field) Nodes() annotation.Node {
	return nil
}

var _ Annotated = (*Field)(nil)

func (f *Field) CheckFieldType() bool {
	field := f.Raw.(*ast.Field)
	t := reflect.TypeOf(field.Type)

	if t.Kind() == reflect.Ptr {
		t = t.Elem() // 如果是指针，获取其基础类型
	}

	if t.Kind() == reflect.Struct {
		return true
	} else {
		return false
	}
}
