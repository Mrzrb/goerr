package sub

import (
	"fmt"

	"github.com/Mrzrb/goerr/annotations/aop"
)

// @Aop(type="aspect")
type Demo struct{}

// @Aop(type="around")
func (r *Demo) Handle(joint aop.Jointcut) {
	fmt.Println("before run")
	joint.Fn()
	fmt.Println("after run")
}

// @Aop(type="point", target="Common")
type BisClient struct{}

// @Aop(type="pointcut")
func (b *BisClient) Hello() (int64, error) {
	fmt.Println("this is in func")
	return 655, nil
}
