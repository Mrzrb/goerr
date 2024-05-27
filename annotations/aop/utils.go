package aop

func generateChain(raw func() error, effect func(Jointcut, MuteableArgs) error, joint Jointcut, m MuteableArgs) func() error {
	joint.Fn = raw
	return func() error {
		return effect(joint, m)
	}
}

func GenerateChain(joint *Jointcut, m MuteableArgs, chain ...func(Jointcut, MuteableArgs) error) {
	for _, v := range chain {
		joint.Fn = generateChain(joint.Fn, v, *joint, m)
	}
}
