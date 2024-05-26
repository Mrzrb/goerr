package subdir

import "github.com/Mrzrb/goerr/annotations/aop"

// @Aop(type="aspect")
type SubAspect struct{}

// @Aop(type = "around")
func (b *SubAspect) Handle(joint aop.Jointcut) {
	joint.Fn()
}
