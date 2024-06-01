package main

import (
	_ "github.com/Mrzrb/goerr/internal/processors/autowire"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
