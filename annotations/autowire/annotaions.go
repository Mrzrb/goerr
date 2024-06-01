package autowire

type AutowireMete struct {
	File  string `annotations:"name=val"`
	Scope string `annotations:"name=scope"`
}

type Component struct {
	// Component Name
	Name string `annotations:"name=name"`
}

type Autowired struct {
	// Inject Component's name
	Name string `annotations:"name=name"`
}

type Factory struct {
	Name string `annotations:"name=name"`
}
