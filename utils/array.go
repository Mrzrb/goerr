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

func Walk[T any](src []T, fn func(T)) {
	for _, v := range src {
		fn(v)
	}
}

func Uniq[T any, K comparable](src []T, fn func(T) K) []T {
	var ret []T
	um := map[K]any{}

	for _, v := range src {
		k := fn(v)
		if _, ok := um[k]; !ok {
			um[k] = 1
			ret = append(ret, v)
		}
		continue
	}

	return ret
}

func Filter[T any](arr []T, fn func(T) bool) []T {
	var ret []T
	for _, v := range arr {
		if fn(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func Contains[T comparable](arr []T, item T) bool {
	return len(Filter(arr, func(t T) bool { return item == t })) > 0
}
