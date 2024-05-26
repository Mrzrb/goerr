package basic

import "github.com/Mrzrb/goerr/annotations/aop"

// @Aop(type="aspect")
type BasicAspect struct{}

// @Aop(type = "around")
func (b *BasicAspect) Handle(joint aop.Jointcut) {
	joint.Fn()
}
