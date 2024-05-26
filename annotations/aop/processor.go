package aop

import (
	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

var (
	e         = core.NewExporter(&core.Assembler)
	processor = NewProcessor(&e)
)

func init() {
	annotation.Register[Aop](&processor)
}

// @Constructor
type Processor struct {
	core.Collector // @Init
	core.FileExporter
}

// Name implements annotation.AnnotationProcessor.
func (p *Processor) Name() string {
	return "aop"
}

func getAnnotated[T any](collector *core.Collector, filter ...func(i T) bool) []core.Annotated {
	ret := []core.Annotated{}
	collector.Walk(func(node core.Annotated) {
		if len(annotation.FindAnnotations[T](node.Annotate())) > 0 {
			if len(filter) > 0 {
				annotations := annotation.FindAnnotations[T](node.Annotate())
				for _, v := range filter {
					for _, a := range annotations {
						if !v(a) {
							return
						}
					}
				}
			}
			ret = append(ret, node)
		}
	})
	return ret
}

func (p *Processor) getUnit() []*Unit {
	aspect := utils.Map(getAnnotated[Aop](&p.Collector, func(i Aop) bool {
		return i.Type == "aspect"
	}), func(t core.Annotated) *core.Struct {
		return utils.MustBool(core.Cast[*core.Struct](t))
	})
	point := utils.Map(getAnnotated[Aop](&p.Collector, func(i Aop) bool {
		return i.Type == "point"
	}), func(t core.Annotated) *core.Struct {
		return utils.MustBool(core.Cast[*core.Struct](t))
	})
	pointCut := utils.Map(getAnnotated[Aop](&p.Collector, func(i Aop) bool {
		return i.Type == "pointcut"
	}), func(t core.Annotated) *core.Method {
		return utils.MustBool(core.Cast[*core.Method](t))
	})

	affects := utils.Map(getAnnotated[Aop](&p.Collector, func(i Aop) bool {
		return utils.Contains([]string{string(Before), string(Around), string(After), string(CatchPanic)}, i.Type)
	}), func(t core.Annotated) *core.Method {
		return utils.MustBool(core.Cast[*core.Method](t))
	})

	us := []*Unit{}
	for _, v := range point {
		annos := annotation.FindAnnotations[Aop](v.Annotations())
		if len(annos) > 1 {
			panic("num of aop point aop must be 1")
		}
		anno := annos[0]
		annoStruct := core.GetByName(utils.Map(aspect, func(t *core.Struct) core.Annotated { return t }), anno.Target).(*core.Struct)
		u := &Unit{
			BaseOutputer: core.BaseOutputer{
				File:    v.DstFileName("aop"),
				Imports: v.Imports(),
				Package: v.Meta().PackageName(),
			},
			Struct:     *v,
			Aspect:     *annoStruct,
			AspectType: annoStruct.Name,
			Method:     []core.BaseFuncOutputer{},
		}
		methods := utils.Map(core.GetMethod(utils.Map(pointCut, func(t *core.Method) core.Annotated { return t }), "*"+v.Name), func(t core.Annotated) *core.Method { return t.(*core.Method) })
		for _, m := range methods {
			u.Method = append(u.Method, core.BaseFuncOutputer{
				IsMethod:  true,
				Receiver:  m.Receiver.Name,
				FuncName:  m.FuncIdent.Name,
				Params:    m.Param,
				Param:     "",
				Returns:   m.Retern,
				Return:    "",
				CallParam: "",
				HasReturn: false,
			})
			u.BaseOutputer.Imports = append(u.BaseOutputer.Imports, m.Imports()...)
		}
		if len(methods) > 0 {
			affect := utils.Filter(affects, func(m *core.Method) bool {
				return m.Receiver.Type == "*"+u.AspectType
			})
			if len(affect) > 1 || len(affect) == 0 {
				panic("affect must be 1")
			}
			u.Affect = *affect[0]
			u.AspectType = affects[0].Name
			if methods[0].Meta().PackageName() != annoStruct.Meta().PackageName() {
				// u.BaseOutputer.Imports = append(u.BaseOutputer.Imports, annoStruct.Meta().PackageName())
				u.PushImport(annoStruct.Meta().PackageName(), annoStruct.Meta().Dir())
			}
		}

		us = append(us, u)
	}
	return us
}

// Output implements annotation.AnnotationProcessor.
func (p *Processor) Output() map[string][]byte {
	for _, v := range p.getUnit() {
		p.FileExporter.Append(v)
	}
	return p.Export()
}

// Process implements annotation.AnnotationProcessor.
func (p *Processor) Process(node annotation.Node) error {
	if len(annotation.FindAnnotations[Aop](node.Annotations())) <= 0 {
		return nil
	}
	p.Collector.Process(node)
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (p *Processor) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*Processor)(nil)
