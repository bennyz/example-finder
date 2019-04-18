package util

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	arr1 := []int64{1, 2, 3, 4, 5}
	arr2 := []int64{4, 5, 7}
	got := Difference(arr1, arr2)
	expected := []int64{7}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Difference(arr1, arr2) = %v, want %v", got, expected)
	}
}
