package autowire

import (
	"github.com/Mrzrb/goerr/annotations/autowire"
	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

var process = Process{
	Collector: core.Collector{},
}

type Process struct {
	Collector core.Collector
	Pre       PreCollector
	Assembler Assembler
}

// Name implements annotation.AnnotationProcessor.
func (p *Process) Name() string {
	return "autowire"
}

// Output implements annotation.AnnotationProcessor.
func (p *Process) Output() map[string][]byte {
	p.Prepare()
	return nil
}

// Process implements annotation.AnnotationProcessor.
func (p *Process) Process(node annotation.Node) error {
	p.Collector.Process(node)
	if len(annotation.FindAnnotations[autowire.AutowireMete](node.Annotations())) > 0 {
		anno := annotation.FindAnnotations[autowire.AutowireMete](node.Annotations())[0]
		p.Assembler.File = anno.File
		p.Assembler.Scope = anno.Scope
	}
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (p *Process) Version() string {
	return "0.1"
}

// collect scan things need to generate
func (p *Process) Prepare() {
	components := []*Struct{}
	factories := []*core.Func{}
	fields := []*core.Field{}
	p.Collector.Walk(func(node core.Annotated) {
		if s, ok := core.Cast[*core.Struct](node); ok {
			p.Pre.Components = append(p.Pre.Components, s.Id())
			components = append(components, &Struct{
				Struct: *s,
				Depend: DependencyTree{},
			})
		}
		if s, ok := core.Cast[*core.Func](node); ok {
			p.Pre.Factories = append(p.Pre.Components, s.Id())
			factories = append(factories, s)
		}
		if s, ok := core.Cast[*core.Field](node); ok {
			p.Pre.Factories = append(p.Pre.Components, s.Id())
			fields = append(fields, s)
		}
	})

	finder := func(src []*Struct, id string) *Struct {
		find := utils.Filter(src, func(s *Struct) bool {
			return id == s.Id()
		})
		if len(find) == 0 {
			return nil
		}
		return find[0]
	}

	// buildDenpendency
	p.Assembler.Components = components
	p.Assembler.Factory = factories
	for _, v := range p.Assembler.Components {
		v.WalkFieldAnnoted(func(field core.Field) {
			if !field.CheckFieldType() {
				return
			}
			depend := finder(p.Assembler.Components, field.Id())
			v.Depend.ChildDependencies = append(v.Depend.ChildDependencies, &DependencyTree{
				Id:                depend.Id(),
				ChildDependencies: []*DependencyTree{},
			})
		})
	}
}

func (p *Process) GetComponentById(id string) *Struct {
	for _, c := range p.Assembler.Components {
		if c.Id() == id {
			return c
		}
	}
	return nil
}

var _ annotation.AnnotationProcessor = (*Process)(nil)
