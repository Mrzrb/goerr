package utils

func InsertInPos[T any](arr *[]T, s T, pos int64) {
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
