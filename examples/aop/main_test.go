package main

import (
	"testing"
)

func TestProxy(t *testing.T) {
	d := NewBisClientProxy(&BisClient{})
	Run(d)
}

func Run(d BisClientInterface) {
	d.Hello()
	d.Hello2()
}
