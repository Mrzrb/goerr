package twoaspect

import (
	"fmt"

	"github.com/Mrzrb/goerr/examples/aop/common"
)

type Two1 struct{}

func (b *Two1) Hello() (int64, error) {
	fmt.Println("this is in func")
	return 655, nil
}

// @Aop(type="point", target="Logger")
type Two2 struct {
	This common.Common
}

// @Aop(type="pointcut")
func (b *Two2) Hello(param1 int, s1 *Two1) (int64, error) {
	fmt.Printf("param is param1: %d, s1: %+v", param1, s1)
	param1 += 4
	return int64(param1), nil
}
