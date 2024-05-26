package twoaspect

import "testing"

func TestTwo(t *testing.T) {
	g := NewTwo2Proxy(&Two2{})
	g.Hello(123, Two1{})
}
