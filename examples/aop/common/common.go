package common

import (
	"fmt"

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
func (c *Logger) Handler(joint aop.Jointcut, param *aop.RunContext) (err error) {
	param.MuteableArgs.Args[0].Value = 55
	err = joint.Fn()
	fmt.Printf("\nresult: %+v %+v\n", param.ReturnResult.Args[0], param.ReturnResult.Args[1])
	if err != nil {
		panic(err)
	}
	return err
}
