package aop

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
