package twoaspect

// @Aop(type="point", target="BaseAspect")
type Service struct{}

// @Aop(type="pointcut")
func (s *Service) Hello() error {
	return nil
}

// @Aop(type="point", target="BaseAspect1")
type Service1 struct{}

// @Aop(type="pointcut")
func (s *Service) Hello1() error {
	return nil
}
