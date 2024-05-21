package main

import (
	_ "github.com/Mrzrb/goerr/annotations/getset"
	_ "github.com/Mrzrb/goerr/annotations/logger"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func main() {
	annotation.Process()
}
