package getset

import (
	"errors"
	"fmt"
	"go/ast"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

func init() {
	processor := NewGsProcessor(&oIns)
	annotation.Register[Getter](processor)
	annotation.Register[Setter](processor)
	annotation.Register[GetterSetter](processor)
}

type GsProcessor struct {
	Parsed  []core.Annotated
	Targets []gsTarget
	Codes   map[string]string
	Out     *out
	core.Collector
}

func NewGsProcessor(o *out) *GsProcessor {
	returnValue := GsProcessor{
		Parsed:    []core.Annotated{},
		Targets:   []gsTarget{},
		Codes:     map[string]string{},
		Out:       o,
		Collector: core.Collector{},
	}
	o.g = &returnValue

	return &returnValue
}

type gsTarget struct {
	field []struct {
		Name       string
		Type       string
		IsGetter   bool
		IsSetter   bool
		TargetName string
	}
	packageName string
	file        string
	dir         string
	Imps        utils.DistinctImports
}

func (g *gsTarget) FullFileName() string {
	return g.dir + "/" + g.file
}

// Name implements annotation.AnnotationProcessor.
func (g *GsProcessor) Name() string {
	return "GetterSetter"
}

// Output implements annotation.AnnotationProcessor.
func (g *GsProcessor) Output() map[string][]byte {
	ret := g.Out.Out()
	return ret
}

// Process implements annotation.AnnotationProcessor.
func (g *GsProcessor) Process(node annotation.Node) error {
	g.Collector.Process(node)
	if v := core.Parse(node); v != nil {
		g.Parsed = append(g.Parsed, v)
	}
	// s := core.NewStruct(node)
	// fmt.Println(s)
	// e1 := g.gsProcess(node)
	// e2 := g.gProcess(node)
	// e3 := g.sProcess(node)
	// return utils.Or(
	// 	e1, e2, e3,
	// )
	//
	return nil
}

func preProcess[T any](node annotation.Node, g *GsProcessor) error {
	ans := annotation.FindAnnotations[T](node.Annotations())
	if len(ans) == 0 {
		return nil
	}
	if len(ans) > 1 {
		return errors.New("annotation > 0 is not permit")
	}
	ns, ok := annotation.CastNode[*ast.TypeSpec](node)
	if !ok {
		return nil
	}
	nss := ns.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	var v any
	v = Getter{}
	_, ok = (v).(T)
	if ok {
		g.appendStruct(node, nss, true, false)
	}

	v = Setter{}
	_, ok = (v).(T)
	if ok {
		g.appendStruct(node, nss, false, true)
	}

	v = GetterSetter{}
	_, ok = (v).(T)
	if ok {
		g.appendStruct(node, nss, true, true)
	}

	return nil
}

func extractTypeFromExpr(v ast.Expr) string {
	switch t := v.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return extractTypeFromExpr(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + extractTypeFromExpr(t.X)
	case *ast.ArrayType:
		return "[]" + extractTypeFromExpr(t.Elt)
	default:
		return fmt.Sprintf("%T", t)
	}
}

func (g *GsProcessor) appendStruct(node annotation.Node, nss *ast.StructType, getter bool, setter bool) error {
	gs := gsTarget{
		field: []struct {
			Name       string
			Type       string
			IsGetter   bool
			IsSetter   bool
			TargetName string
		}{},
		packageName: "",
		file:        "",
		dir:         "",
		Imps:        map[utils.Import]struct{}{},
	}
	for _, v := range nss.Fields.List {
		n := node.ASTNode().(*ast.TypeSpec)
		gs.Imps.Merge(utils.GetImports(v.Type, node.Lookup().FindImportByAlias))
		gs.field = append(gs.field, struct {
			Name       string
			Type       string
			IsGetter   bool
			IsSetter   bool
			TargetName string
		}{
			Name:       v.Names[0].Name,
			Type:       extractTypeFromExpr(v.Type),
			IsGetter:   getter,
			IsSetter:   setter,
			TargetName: n.Name.Name,
		})
		gs.file = node.Meta().FileName()
		gs.dir = node.Meta().Dir()
		gs.packageName = node.Meta().PackageName()
	}
	g.Targets = append(g.Targets, gs)
	return nil
}

func (g *GsProcessor) gsProcess(node annotation.Node) error {
	return preProcess[GetterSetter](node, g)
}

func (g *GsProcessor) gProcess(node annotation.Node) error {
	return preProcess[Getter](node, g)
}

func (g *GsProcessor) sProcess(node annotation.Node) error {
	return preProcess[Setter](node, g)
}

func (g *GsProcessor) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*GsProcessor)(nil)
