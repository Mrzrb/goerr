package common

import (
	"fmt"
	"time"

	"github.com/Mrzrb/goerr/annotations/aop"
)

// @Aop(type="aspect")
type Common struct{}

// @Aop(type="around")
func (c *Common) Handler(joint aop.Jointcut) error {
	var err error
	fmt.Println("coomon enter")
	err = joint.Fn()
	fmt.Println("coomon exit")
	return err
}

// @Aop(type="aspect")
type Common1 struct{}

// @Aop(type="around")
func (c *Common1) Handler(joint aop.Jointcut) (err error) {
	fmt.Println("coomon1 enter")
	err = joint.Fn()
	fmt.Println("coomon1 exit")
	return err
}

func GenerateChain(fn func()) func() {
	return func() {
		fn()
	}
}

// @Aop(type="aspect")
type Logger struct{}

// @Aop(type="around")
func (c *Logger) Handler(joint aop.Jointcut) (err error) {
	now := time.Now().Unix()
	fmt.Printf("start exec in %d", now)
	defer func() {
		stop := time.Now().Unix()
		fmt.Printf("end exec in %d, duration %d", stop, stop-now)
	}()
	err = joint.Fn()
	time.Sleep(time.Second)
	if err != nil {
		panic(err)
	}
	return err
}
