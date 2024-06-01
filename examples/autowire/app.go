package main

// @Component()
type App struct {
	// @Autowired()
	Name Name
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
