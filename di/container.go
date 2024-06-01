package di

import "github.com/samber/do/v2"

var injectMap = map[string]*do.Scope{}

var GlobalInjector = do.New()

func GetInject(name string) do.Injector {
	if i, ok := injectMap[name]; ok {
		return i
	}
	scoped := GlobalInjector.Scope(name)
	injectMap[name] = scoped

	return scoped
}
