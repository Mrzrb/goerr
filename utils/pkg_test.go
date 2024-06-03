package utils

import (
	"fmt"
	"testing"
)

func TestPkg(t *testing.T) {
	pkg := GetFullPackage(".")
	fmt.Println(pkg)
}
