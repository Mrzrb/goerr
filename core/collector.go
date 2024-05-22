package core

import annotation "github.com/YReshetko/go-annotation/pkg"

type Collector struct {
	Annotations []Annotated
}

// Process implements annotation.AnnotationProcessor.
func (c *Collector) Process(node annotation.Node) error {
	if v := Parse(node); v != nil {
		c.Annotations = append(c.Annotations, v)
	}
	return nil
}

func (c *Collector) Walk(fn func(node Annotated)) {
	for _, n := range c.Annotations {
		fn(n)
	}
}
