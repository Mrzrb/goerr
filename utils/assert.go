package utils

func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func Or(errs ...error) error {
	for _, v := range errs {
		if v != nil {
			return v
		}
	}
	return nil
}
