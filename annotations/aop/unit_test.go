package aop

import (
	"fmt"
	"testing"

	"github.com/Mrzrb/goerr/core"
)

func TestOutputUnit(t *testing.T) {
	u := &Unit{
		BaseOutputer: core.BaseOutputer{
			File:    "",
			Imports: []string{},
			Package: "",
		},
		Struct: core.Struct{
			Node: core.Node{},
			Ident: core.Ident{
				AnnotationsMix: core.AnnotationsMix{},
				Name:           "Demo",
				Type:           "DemoType",
				Raw:            nil,
			},
			Field: []core.Field{},
		},
		Aspect: core.Struct{
			Node: core.Node{},
			Ident: core.Ident{
				AnnotationsMix: core.AnnotationsMix{},
				Name:           "DemoAspect",
				Type:           "DemoAspectType",
				Raw:            nil,
			},
			Field: []core.Field{},
		},
		Method: []core.BaseFuncOutputer{
			{
				IsMethod:  false,
				Receiver:  "",
				FuncName:  "Hello",
				Params:    []core.Ident{},
				Param:     "",
				Returns:   []core.Ident{},
				Return:    "",
				CallParam: "",
				HasReturn: false,
			},
			{
				IsMethod:  false,
				Receiver:  "",
				FuncName:  "Hello2",
				Params:    []core.Ident{},
				Param:     "",
				Returns:   []core.Ident{},
				Return:    "",
				CallParam: "",
				HasReturn: false,
			},
		},
	}

	fmt.Println(string(u.Output()))
}
