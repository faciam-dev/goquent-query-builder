package sliceutils

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func RemoveIfContains[T comparable](elems []T, v T) []T {
	for i := 0; i < len(elems); i++ {
		if elems[i] == v {
			elems = append(elems[:i], elems[i+1:]...)
			i--
		}
	}

	return elems
}

func Reverse[T comparable](s []T) []T {
	length := len(s)
	reversed := make([]T, length)

	for i, v := range s {
		reversed[length-1-i] = v
	}

	return reversed
}
