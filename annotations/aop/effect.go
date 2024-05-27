package aop

import "github.com/Mrzrb/goerr/annotations/aop_core"

type Effecter interface {
	Around(joint aop_core.Jointcut)
	Before(joint aop_core.Jointcut)
	After(joint aop_core.Jointcut)
	Catch(joint aop_core.Jointcut)
}
