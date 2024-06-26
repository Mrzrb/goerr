package main

import (
	_ "github.com/Mrzrb/goerr/annotations/aop"
	_ "github.com/Mrzrb/goerr/annotations/getset"
	_ "github.com/Mrzrb/goerr/annotations/logger"
	_ "github.com/Mrzrb/goerr/internal/processors/autowire"
	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	_ "github.com/YReshetko/go-annotation/annotations/mapper"
	_ "github.com/YReshetko/go-annotation/annotations/validator"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
