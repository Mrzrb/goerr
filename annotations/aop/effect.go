package aop

type Effecter interface {
	Around(joint Jointcut)
	Before(joint Jointcut)
	After(joint Jointcut)
	Catch(joint Jointcut)
}
