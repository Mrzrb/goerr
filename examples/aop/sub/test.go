package sub

import (
	"fmt"
)

// @Aop(type="point", target="Common")
type BisClient struct{}

// @Aop(type="pointcut")
func (b *BisClient) Hello() (int64, error) {
	fmt.Println("this is in func")
	return 655, nil
}
