package common

import (
	"fmt"

	"github.com/Mrzrb/goerr/annotations/aop"
)

// @Aop(type="aspect")
type Common struct{}

// @Aop(type="around")
func (c *Common) Handler(joint aop.Jointcut) {
	joint.Fn()
}

// @Aop(type="aspect")
type Common1 struct{}

// @Aop(type="around")
func (c *Common1) Handler(joint aop.Jointcut) {
	fmt.Printf("this is common 1, type %s , args :%+v\n", joint.TargetType, joint.Args)
	joint.Fn()
	fmt.Println("this is common 2")
}
