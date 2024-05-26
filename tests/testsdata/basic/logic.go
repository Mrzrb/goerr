package basic

// @Aop(type="point", target="BasicAspect")
type Service struct{}

// @Aop(type="pointcut")
func (s *Service) Hello() error {
	return nil
}
