package autowire

type AutowireMete struct {
	File    string `annotation:"name=file"`
	Scope   string `annotation:"name=scope"`
	Package string `annotation:"name=package"`
}

type Component struct {
	// Component Name
	Name string `annotation:"name=name"`
}

type Autowired struct {
	// Inject Component's name
	Name string `annotation:"name=name"`
}

type Factory struct {
	Name string `annotation:"name=name"`
}
