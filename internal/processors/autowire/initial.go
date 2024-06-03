package autowire

import (
	anno "github.com/Mrzrb/goerr/annotations/autowire"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	annotation.Register[anno.AutowireMete](&process)
	annotation.Register[anno.Factory](&process)
	annotation.Register[anno.Autowired](&process)
	annotation.Register[anno.Component](&process)
}
