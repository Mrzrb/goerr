package core

var (
	_ Identity = (*Struct)(nil)
	_ Identity = (*Method)(nil)
	_ Identity = (*Func)(nil)
	_ Identity = (*Field)(nil)

	_ Callable = (*Method)(nil)
	_ Callable = (*Func)(nil)
)
