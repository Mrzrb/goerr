package main

import (
	_ "github.com/Mrzrb/goerr/annotations/getset"
	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
