package aop

func generateChain(raw func() error, effect func(Jointcut) error, joint Jointcut) func() error {
	joint.Fn = raw

	return func() error {
		return effect(joint)
	}
}

func GenerateChain(joint Jointcut, chain ...func(Jointcut) error) func() error {
	fn := func() error {
		return joint.Fn()
	}

	for _, v := range chain {
		fn = generateChain(fn, v, joint)
	}

	return fn
}
