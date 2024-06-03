package autowire

import (
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/annotations/autowire"
	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

var e = core.NewExporter(&core.Assembler)

var process = Process{
	Collector: core.Collector{},
	Pre:       PreCollector{},
	Assembler: &Assembler{
		Scope:        "",
		Pre:          PreCollector{},
		Components:   []*Struct{},
		Factory:      []*core.Func{},
		BaseOutputer: core.BaseOutputer{},
		b:            strings.Builder{},
		inited:       map[string]initItem{},
	},
	FileExporter: &e,
}

type Process struct {
	Collector core.Collector
	Pre       PreCollector
	Assembler *Assembler
	core.FileExporter
}

// Name implements annotation.AnnotationProcessor.
func (p *Process) Name() string {
	return "autowire"
}

// Output implements annotation.AnnotationProcessor.
func (p *Process) Output() map[string][]byte {
	p.Prepare()
	p.Append(p.Assembler)
	return p.Export()
}

// Process implements annotation.AnnotationProcessor.
func (p *Process) Process(node annotation.Node) error {
	p.Collector.Process(node)
	if len(annotation.FindAnnotations[autowire.AutowireMete](node.Annotations())) > 0 {
		anno := annotation.FindAnnotations[autowire.AutowireMete](node.Annotations())[0]
		p.Assembler.BaseOutputer.File = anno.File
		p.Assembler.BaseOutputer.Package = anno.Package
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
	fields := map[string]*core.Field{}
	p.Collector.Walk(func(node core.Annotated) {
		if s, ok := core.Cast[*core.Struct](node); ok {
			p.Pre.Components = append(p.Pre.Components, s.Id())
			components = append(components, &Struct{
				Struct: *s,
				Depend: DependencyTree{},
			})
			for _, f := range s.Field {
				if field, ok := f.Raw.(*ast.Field); ok {
					if field.Doc == nil {
						continue
					}
					for _, comment := range field.Doc.List {
						if strings.Contains(comment.Text, "Autowire") {
							ff := &core.Field{
								Ident: core.Ident{
									AnnotationsMix: core.AnnotationsMix{},
									Name:           f.Name,
									Type:           f.Type,
									Raw:            field,
									IsPointer:      f.IsPointer,
									Package:        "",
								},
								Package:         node.Nodes().Meta().PackageName(),
								Alias:           "",
								FullPackagePath: "",
								Parent:          s.Node,
							}
							fields[strings.ReplaceAll(f.Type, "*", "")] = ff
						}
					}
				}
			}
			return
		}
		if s, ok := core.Cast[*core.Func](node); ok {
			p.Pre.Factories = append(p.Pre.Components, s.Id())
			factories = append(factories, s)
			return
		}
		if s, ok := core.Cast[*core.Field](node); ok {
			fields[s.Id()] = s
			return
		}
	})

	finder := func(src []*Struct, id string) *Struct {
		find := utils.Filter(src, func(s *Struct) bool {
			sid := s.Id()
			return id == sid
		})
		if len(find) == 0 {
			return nil
		}
		return find[0]
	}

	// buildDenpendency
	p.Assembler.Components = components
	p.Assembler.Factory = factories
	p.Assembler.Fields = fields
	for _, v := range p.Assembler.Components {
		v.WalkFieldAnnoted(func(field core.Field) {
			if !field.CheckFieldType() {
				return
			}
			depend := finder(p.Assembler.Components, strings.ReplaceAll(field.Type, "*", ""))
			if depend == nil {
				return
			}
			v.Depend.ChildDependencies = append(v.Depend.ChildDependencies, &DependencyTree{
				Id:                depend.Id(),
				Name:              field.Name,
				ChildDependencies: []*DependencyTree{},
			})
		})
	}
}

var _ annotation.AnnotationProcessor = (*Process)(nil)
