package util

import (
	"strconv"
	"strings"
)

func StringToIntSlice(s string) ([]int64, error) {
	parts := strings.Split(s, ",")

	var result []int64

	for _, part := range parts {
		num, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, num)
	}

	return result, nil
}
