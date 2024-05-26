// Code generated by aop annotation processor. DO NOT EDIT.
// versions:
//		go: go1.21.0
//		go-annotation: 0.1.0
//		aop: 0.1

package service

import (
	"github.com/Mrzrb/goerr/annotations/aop"
	"github.com/Mrzrb/goerr/tests/testsdata/subdir"
)

type ServiceProxy struct {
	inner  *Service
	aspect *subdir.SubAspect
}

func NewServiceProxy(inner *Service) *ServiceProxy {
	return &ServiceProxy{
		inner:  inner,
		aspect: &subdir.SubAspect{},
	}
}

type ServiceInterface interface {
	Hello() (ret1 error)
}

func (r *ServiceProxy) Hello() (ret1 error) {
	joint := aop.Jointcut{
		TargetName: "Service",
		TargetType: "Service",
		Args:       []aop.Args{},
		Fn: func() {
			ret1 = r.inner.Hello()
		},
	}

	r.aspect.Handle(joint)

	return ret1
}
