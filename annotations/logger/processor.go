package logger

import (
	"github.com/Mrzrb/goerr/core"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

// @Constructor
type LoggerProcess struct {
	core.Collector // @Init
	core.FileExporter
}

func init() {
	a := core.NewAssembler()
	e := core.NewExporter(&a)
	p := NewLoggerProcess(&e)
	annotation.Register[Logger](&p)
}

// Name implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Name() string {
	return "logger"
}

// Output implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Output() map[string][]byte {
	ret := l.FileExporter.Export()
	return ret
}

// Process implements annotation.AnnotationProcessor.
// Subtle: this method shadows the method (Collector).Process of LoggerProcess.Collector.
func (l *LoggerProcess) Process(node annotation.Node) error {
	l.Collector.Process(node)
	u := Unit{}
	l.FileExporter.Append(&u)
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*LoggerProcess)(nil)
