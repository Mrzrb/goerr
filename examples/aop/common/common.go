package common

import (
	"fmt"
	"time"

	"github.com/Mrzrb/goerr/annotations/aop"
)

// @Aop(type="aspect")
type Common struct{}

// @Aop(type="around")
func (c *Common) Handler(joint aop.Jointcut) {
	fmt.Println("coomon enter")
	joint.Fn()
	fmt.Println("coomon exit")
}

// @Aop(type="aspect")
type Common1 struct{}

// @Aop(type="around")
func (c *Common1) Handler(joint aop.Jointcut) {
	fmt.Println("coomon1 enter")
	joint.Fn()
	fmt.Println("coomon1 exit")
}

func GenerateChain(fn func()) func() {
	return func() {
		fn()
	}
}

// @Aop(type="aspect")
type Logger struct{}

// @Aop(type="around")
func (c *Logger) Handler(joint aop.Jointcut) {
	now := time.Now().Unix()
	fmt.Printf("start exec in %d", now)
	defer func() {
		stop := time.Now().Unix()
		fmt.Printf("end exec in %d, duration %d", stop, stop-now)
	}()
	joint.Fn()
	time.Sleep(time.Second)
}
