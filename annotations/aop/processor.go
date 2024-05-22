package aop

import (
	"bytes"
	"fmt"

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

// Output implements annotation.AnnotationProcessor.
func (p *Processor) Output() map[string][]byte {
	return nil
}

// Process implements annotation.AnnotationProcessor.
func (p *Processor) Process(node annotation.Node) error {
	b := bytes.NewBufferString("")
	utils.AstToGo(b, node.ASTNode())
	fmt.Println(b.String())
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
