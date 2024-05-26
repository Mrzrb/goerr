package twoaspect

import (
	"fmt"
)

type Two1 struct{}

func (b *Two1) Hello() (int64, error) {
	fmt.Println("this is in func")
	return 655, nil
}

// @Aop(type="point", target="Common,Common1")
type Two2 struct{}

// @Aop(type="pointcut")
func (b *Two2) Hello(param1 int, s1 Two1) (int64, error) {
	fmt.Println("this is in func")
	return 655, nil
}
