package testsdata

import (
	"github.com/Mrzrb/goerr/annotations/aop_core"
)

// @Aop(type="aspect")
type BaseAspect struct{}

// @Aop(type = "around")
func (b *BaseAspect) Handle(joint aop_core.Jointcut) {
	joint.Fn()
}

// @Aop(type="aspect")
type Basespect1 struct{}

// @Aop(type = "around")
func (b *Basespect1) Handle1(joint aop_core.Jointcut) {
	joint.Fn()
}
