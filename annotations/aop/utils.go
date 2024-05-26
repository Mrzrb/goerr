package aop

func generateChain(raw func(), effect func(Jointcut), joint Jointcut) func() {
	joint.Fn = raw

	return func() {
		effect(joint)
	}
}

func GenerateChain(joint Jointcut, chain ...func(Jointcut)) func() {
	fn := func() {
		joint.Fn()
	}

	for _, v := range chain {
		fn = generateChain(fn, v, joint)
	}

	return fn
}
