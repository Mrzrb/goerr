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
	joint.Args[0].Value = 555
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

// @Aop(type="aspect")
type Logger struct{}

// @Aop(type="around")
func (c *Logger) Handler(joint aop.Jointcut, param aop.MuteableArgs) (err error) {
	now := time.Now().Unix()
	param.Args[0].Value = 55
	fmt.Printf("start exec in %d", now)
	defer func() {
		stop := time.Now().Unix()
		fmt.Printf("end exec in %d, duration %d", stop, stop-now)
	}()
	err = joint.Fn()
	if err != nil {
		panic(err)
	}
	return err
}
