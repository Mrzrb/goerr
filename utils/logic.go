package utils

func OrGet[T any](is bool, a T, b T) T {
	if is {
		return a
	}
	return b
}
