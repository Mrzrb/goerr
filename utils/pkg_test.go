package utils

import (
	"fmt"
	"testing"
)

func TestPkg(t *testing.T) {
	pkg := GetFullPackage("../examples/aop/sub")
	fmt.Println(pkg)
}
