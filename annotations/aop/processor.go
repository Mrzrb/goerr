package aop

import (
	"strings"

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
		annoStructs := utils.Filter(aspect, func(s *core.Struct) bool {
			return utils.Contains(strings.Split(anno.Target, ","), s.Name)
		})
		methods := utils.Map(core.GetMethod(utils.Map(pointCut, func(t *core.Method) core.Annotated { return t }), "*"+v.Name), func(t core.Annotated) *core.Method { return t.(*core.Method) })
		u := &Unit{
			BaseOutputer: core.BaseOutputer{
				File:    v.DstFileName("aop"),
				Imports: v.Imports(),
				Package: v.Meta().PackageName(),
			},
			Struct: *v,
			Method: []core.BaseFuncOutputer{},
		}
		utils.Each(annoStructs, func(s *core.Struct) {
			u.Effects = append(u.Effects, Effect{
				Aspect:     *s,
				Affect:     core.Method{},
				AspectType: s.Name,
			})
			if len(methods) > 0 {
				for k, vv := range u.Effects {
					affect := utils.Filter(affects, func(m *core.Method) bool {
						return m.Receiver.Type == "*"+vv.AspectType
					})
					if len(affect) == 0 {
						continue
					}
					vv.Affect = *affect[0]
					vv.AspectType = affects[0].Name
					if methods[0].Meta().PackageName() != s.Meta().PackageName() {
						// u.BaseOutputer.Imports = append(u.BaseOutputer.Imports, annoStruct.Meta().PackageName())
						u.PushImport(s.Meta().PackageName(), s.Meta().Dir())
					}
					u.Effects[k] = vv
				}
			}
		})
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
				Method:    *m,
			})
			u.BaseOutputer.Imports = append(u.BaseOutputer.Imports, m.Imports()...)
		}

		us = append(us, u)
	}
	return us
}

// Output implements annotation.AnnotationProcessor.
func (p *Processor) Output() map[string][]byte {
	for _, v := range p.getUnit() {
		p.Append(v)
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
