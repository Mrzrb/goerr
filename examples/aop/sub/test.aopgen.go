// Code generated by aop annotation processor. DO NOT EDIT.
// versions:
//		go: go1.21.0
//		go-annotation: 0.1.0
//		aop: 0.1

package sub

import (
	"github.com/Mrzrb/goerr/annotations/aop"
	"github.com/Mrzrb/goerr/examples/aop/common"
)

type BisClientProxy struct {
	inner *BisClient

	aspect0 *common.Common
}

func NewBisClientProxy(inner *BisClient) *BisClientProxy {
	return &BisClientProxy{
		inner: inner,

		aspect0: &common.Common{},
	}
}

type BisClientInterface interface {
	Hello() (ret1 int64, ret2 error)
}

func (r *BisClientProxy) Hello() (ret1 int64, ret2 error) {
	joint := aop.Jointcut{
		TargetName: "BisClient",
		TargetType: "BisClient",
		MethodName: "Hello",
		Args:       []aop.Args{},
	}

	mutableArgs := aop.MuteableArgs{}

	joint.Fn = func() error {
		ret1, ret2 = r.inner.Hello()

		if "error" == "error" {
			return ret2
		}
		return nil
	}

	aop.GenerateChain(&joint, mutableArgs,

		func(j aop.Jointcut, m aop.MuteableArgs) error {
			return r.aspect0.Handler(j, m)
		},
	)
	joint.Fn()
	return ret1, ret2
}
