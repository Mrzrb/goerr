package main

import (
	"fmt"

	"github.com/Mrzrb/goerr/annotations/getset"
)

type Ts struct {
	Field1 string
	Field2 getset.Getter
	Test   any
}

// @Logger
func (t *Ts) TestMethod(code int) {
	fmt.Println("tt")
}
