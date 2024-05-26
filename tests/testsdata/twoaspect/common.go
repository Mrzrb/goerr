package twoaspect

import "github.com/Mrzrb/goerr/annotations/aop"

// @Aop(type="aspect")
type BaseAspect struct{}

// @Aop(type = "around")
func (b *BaseAspect) Handle(joint aop.Jointcut) {
	joint.Fn()
}

// @Aop(type="aspect")
type BaseAspect1 struct{}

// @Aop(type = "around")
func (b *BaseAspect1) Handle1(joint aop.Jointcut) {
	joint.Fn()
}
