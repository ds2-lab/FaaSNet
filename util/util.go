package util

func Abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func Max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func IsExist(strings []string, val string) (int, bool) {
	for i, s := range strings {
		if s == val {
			return i, true
		}
	}
	return -1, false
}
