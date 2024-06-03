package main

// @Component()
type App struct {
	// @Autowired()
	Name Name
	// @Autowired()
	Name1 *Name1
}

// @Component()
type Name1 struct {
	FirstName  string
	SecondName string
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
