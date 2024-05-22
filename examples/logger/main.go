package main

import (
	"github.com/Mrzrb/goerr/annotations/getset"
)

type Ts struct {
	Field1 string
	Field2 getset.Getter
	Test   any
}
