package main

import "github.com/Mrzrb/goerr/annotations/aop"

// @Aop(type="aspect")
type Demo struct{}

// @Aop(type="around")
func (r *Demo) Handle(joint aop.Jointcut) {
}

type BisClient struct{}

// @A mix
func (b *BisClient) Hello() error {
	return nil
}

type BasClientProxy struct {
	inner *BisClient
	a     *Demo
}

func (b *BasClientProxy) Hello() {
	joint := aop.Jointcut{
		Fn: func() {
			b.inner.Hello()
		},
	}
	b.a.Handle(joint)
}

type BisClientProxy struct {
	p BisClient
}

func (b *BisClientProxy) Hello() {
}
