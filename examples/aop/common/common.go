package common

import (
	"fmt"

	"github.com/Mrzrb/goerr/annotations/aop"
)

// @Aop(type="aspect")
type Common struct{}

// @Aop(type="around")
func (c *Common) Handler(joint aop.Jointcut) {
	fmt.Println("coomon enter")
	joint.Fn()
	fmt.Println("coomon exit")
}

// @Aop(type="aspect")
type Common1 struct{}

// @Aop(type="around")
func (c *Common1) Handler(joint aop.Jointcut) {
	fmt.Println("coomon1 enter")
	joint.Fn()
	fmt.Println("coomon1 exit")
}

func GenerateChain(fn func()) func() {
	return func() {
		fn()
	}
}
