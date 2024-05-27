package aop_core

func generateChain(raw func() error, effect func(Jointcut, *RunContext) error, joint Jointcut, m *RunContext) func() error {
	joint.Fn = raw
	return func() error {
		return effect(joint, m)
	}
}

func GenerateChain(joint *Jointcut, m *RunContext, chain ...func(Jointcut, *RunContext) error) {
	for _, v := range chain {
		joint.Fn = generateChain(joint.Fn, v, *joint, m)
	}
}

func Cast[T any](src any) T {
	if v, ok := src.(T); ok {
		return v
	}

	var vv T
	return vv
}
