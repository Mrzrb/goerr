# Goal

use annotation to generate code

## Autowire

This is contains some annotation

```go
// This autowireMete annotation set base info that generate code
type AutowireMete struct {
    // file to write
	File    string `annotation:"name=file"`
    // todo
	Scope   string `annotation:"name=scope"`
    // gen code pkg
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
```
