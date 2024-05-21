package main

import "github.com/Mrzrb/goerr/annotations/getset"

type Ts struct {
	Field1 string
	Field2 getset.Getter
	Test   any
	Info   innerStruct
	Infos  *innerStruct
}

type innerStruct struct{}

// @Getter
func TestFunc() (*innerStruct, error) {
	return nil, nil
}
