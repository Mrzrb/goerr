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
}

// Id implements Identity.
func (f *Field) Id() string {
	parentNode, ok := f.Nodes().ParentNode()
	if !ok {
		panic("field must have parent struct")
	}
	pNode := Parse(parentNode)
	if n, ok := Cast[*Struct](pNode); ok {
		return fmt.Sprintf("%s.%s", n.Id(), f.Name)
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
	return f.Nodes()
}

var _ Annotated = (*Field)(nil)

func (f *Field) CheckFieldType() bool {
	field := f.Nodes().ASTNode().(*ast.Field)
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
