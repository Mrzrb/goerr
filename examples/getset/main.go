package main

import "github.com/Mrzrb/goerr/annotations/getset"

// @Getter
type Ts struct {
	Field1 string
	Field2 getset.Getter
	Test   any
}
