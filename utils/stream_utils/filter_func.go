package stream_utils

func Empty[T comparable](s T) bool {
	var v T
	if s == v {
		return true
	}

	return false
}
