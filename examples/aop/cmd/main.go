package main

import (
	_ "github.com/Mrzrb/goerr/annotations/aop"
	_ "github.com/Mrzrb/goerr/annotations/getset"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
