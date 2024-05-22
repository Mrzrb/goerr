package main

// @Aop(type="aspect")
type Demo struct{}

// @Aop(type="before")
func (r *Demo) Handle() {
}

type BisClient struct{}

func (b *BisClient) Hello() {
}

type BisClientProxy struct {
	p BisClient
}

func (b *BisClientProxy) Hello() {
}
