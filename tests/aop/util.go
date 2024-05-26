package aop

import (
	"os/exec"
)

func Run(dir string) (string, error) {
	cmd := exec.Command("go", "run", "../../cmd/main.go", dir)
	result, err := cmd.Output()
	return string(result), err
}
