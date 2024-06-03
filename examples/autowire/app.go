package main

import "github.com/Mrzrb/goerr/examples/autowire/sub"

// @Component()
type App struct {
	// @Autowired()
	Name1 *sub.Sub
}

// @Component()
type Name struct {
	FirstName  string
	SecondName string
}

// @Factory
func NewName() *Name {
	return &Name{}
}
