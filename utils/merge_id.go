package utils

func MergeId(a, b int64) int64 {
	if a > b {
		a, b = b, a
	}

	return a<<10 + b
}
