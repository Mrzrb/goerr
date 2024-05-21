package logger

import (
	"github.com/Mrzrb/goerr/core"
	annotation "github.com/YReshetko/go-annotation/pkg"
)

type LoggerProcess struct {
	core.Collector
}

func init() {
	p := LoggerProcess{}
	annotation.Register[Logger](&p)
}

// Name implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Name() string {
	return "logger"
}

// Output implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Output() map[string][]byte {
	return nil
}

// Process implements annotation.AnnotationProcessor.
// Subtle: this method shadows the method (Collector).Process of LoggerProcess.Collector.
func (l *LoggerProcess) Process(node annotation.Node) error {
	l.Collector.Process(node)
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (l *LoggerProcess) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*LoggerProcess)(nil)
