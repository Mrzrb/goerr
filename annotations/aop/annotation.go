package aop

type Aop struct {
	Type   string `annotation:"name=type,default=aspect,oneOf=aspect;point;pointcut;before;after;around;catchPanic"`
	Target string `annotation:"name=target"`
}
