package service

// @Aop(type="point", target="SubAspect")
type Service struct{}

// @Aop(type="pointcut")
func (s *Service) Hello() error {
	return nil
}
