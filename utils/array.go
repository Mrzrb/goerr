package utils

func InsertInPos[T any](arr *[]T, s T, pos int) {
	if len(*arr) < int(pos)+1 {
		*arr = append(*arr, s)
		return
	}
	var ret []T
	for index, v := range *arr {
		if index == int(pos) {
			ret = append(ret, s, v)
			continue
		}
		ret = append(ret, v)
	}
	*arr = ret
}

func GroupBy[T any, K comparable](arr []T, fn func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, v := range arr {
		k := fn(v)
		result[k] = append(result[k], v)
	}
	return result
}

func Map[T any, R any](src []T, fn func(T) R) []R {
	var ret []R
	for _, v := range src {
		ret = append(ret, fn(v))
	}

	return ret
}
