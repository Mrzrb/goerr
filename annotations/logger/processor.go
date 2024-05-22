package logger

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/Mrzrb/goerr/core"
	"github.com/Mrzrb/goerr/utils"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Constructor
type LoggerProcess struct {
	core.Collector // @Init
	core.FileExporter
}

var (
	a = core.NewAssembler()
	e = core.NewExporter(&a)
	p = NewLoggerProcess(&e)
)

func init() {
	annotation.Register[Logger](&p)
}

// Name implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Name() string {
	return "logger"
}

// Output implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Output() map[string][]byte {
	l.Collector.Walk(func(node core.Annotated) {
		l.Append(parseUnit(node))
	})
	ret := l.Export()
	return ret
}

func parseUnit(node core.Annotated) core.Outputer {
	if _, ok := core.Cast[*core.Struct](node); ok {
		panic("logger can not apply on struct")
	}
	u := Unit{
		IsMethod:  false,
		Receiver:  "",
		FuncName:  "",
		Params:    []core.Ident{},
		Param:     "",
		Returns:   []core.Ident{},
		Return:    "",
		CallParam: "",
	}

	if m, ok := core.Cast[*core.Method](node); ok {
		u.IsMethod = true
		u.Receiver = m.Receiver.Type
		u.FuncName = m.Name
		u.Params = m.Param
		u.Returns = m.Retern
		u.file = m.DstFileName()
	}
	if f, ok := core.Cast[*core.Func](node); ok {
		u.IsMethod = false
		u.FuncName = f.Name
		u.Params = f.Param
		u.Returns = f.Retern
		u.file = f.DstFileName()
	}

	i := utils.DistinctImports{}
	fields := append(u.Params, u.Returns...)
	utils.Walk(fields, func(t core.Ident) {
		f, ok := t.Raw.(*ast.Field)
		if !ok {
			return
		}
		i.Merge(utils.GetImports(f.Type, node.Nodes().Lookup().FindImportByAlias))
		i.Merge(utils.GetImports(f.Type, utils.ImprtTry(node.Nodes())))
	})
	u.Import = utils.Map(i.ToSlice(), func(t utils.Import) string { return t.Package })
	u.HasReturn = len(u.Returns) > 0
	u.Packages = node.Nodes().Meta().PackageName()
	u.Param = strings.Join(utils.Map(u.Params, func(t core.Ident) string {
		return fmt.Sprintf("%s %s", t.Name, t.Type)
	}), ",")
	u.CallParam = strings.Join(utils.Map(u.Params, func(t core.Ident) string {
		return fmt.Sprintf("%s", t.Name)
	}), ",")

	if len(u.Returns) > 0 {
		if len(u.Returns) == 1 {
			u.Return = u.Returns[0].Type
		} else {
			u.Return = fmt.Sprintf("(%s)", strings.Join(utils.Map(u.Returns, func(t core.Ident) string {
				return fmt.Sprintf(" %s", t.Type)
			}), ","))
		}
	}
	return &u
}

// Process implements annotation.AnnotationProcessor.
// Subtle: this method shadows the method (Collector).Process of LoggerProcess.Collector.
func (l *LoggerProcess) Process(node annotation.Node) error {
	fmt.Println(l.FileExporter)
	if _, ok := annotation.CastNode[*ast.FuncDecl](node); !ok {
		panic("logger must applied on func or method")
	}
	l.Collector.Process(node)
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*LoggerProcess)(nil)
