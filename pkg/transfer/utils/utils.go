package utils

import "fmt"

func Min(ints []int64) (int64, error) {
	if len(ints) == 0 {
		return -1, fmt.Errorf("empty input array")
	}

	var result = ints[0]
	for i := 1; i < len(ints); i++ {
		if ints[i] < result {
			result = ints[i]
		}
	}
	return result, nil
}
