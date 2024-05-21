package main

import (
	"fmt"

	"github.com/Mrzrb/goerr/annotations/getset"
)

// @Getter
type Ts struct {
	innerStruct
	Field1 string
	Field2 getset.Getter
	Test   any
	Infos  *innerStruct
}

type innerStruct struct{}

func TestFunc() (*innerStruct, error) {
	return nil, nil
}

func (t *Ts) TestMethod(code int) {
	fmt.Println("tt")
}
