package common

import "github.com/Mrzrb/goerr/annotations/aop"

// @Aop(type="aspect")
type Common struct{}

// @Aop(type="around")
func (c *Common) Handler(joint aop.Jointcut) {
	joint.Fn()
}
