package util

import (
	"strconv"
	"strings"

	backends "github.com/bennyz/example-finder/backend"
)

// Difference returns the difference between two arrays
func Difference(arr1, arr2 []int64) (result []int64) {
	hash := make(map[int64]struct{})

	for _, item := range arr1 {
		hash[item] = struct{}{}
	}

	for _, item := range arr2 {
		if _, ok := hash[item]; !ok {
			result = append(result, item)
		}
	}

	return result
}

// SliceToString converts an int64 array to a comma separated string
func SliceToString(arr []int64) string {
	temp := make([]string, len(arr))
	for i, e := range arr {
		temp[i] = strconv.FormatInt(e, 10)
	}

	return strings.Join(temp, ",")
}

// MapToSlice converts map's values to a slice
func MapToSlice(m map[int64]*backends.Result) []*backends.Result {
	s := make([]*backends.Result, 0, len(m))

	for _, v := range m {
		s = append(s, v)
	}

	return s
}
