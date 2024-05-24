package utils

func Must[T any](t T, e error) T {
	if e != nil {
		panic(e)
	}
	return t
}

func MustBool[T any](t T, b bool) T {
	if b != true {
		panic(b)
	}
	return t
}

func MustPointer[T any](t *T) *T {
	if t == nil {
		panic(t)
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
