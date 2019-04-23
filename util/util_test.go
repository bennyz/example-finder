package util

import (
	"reflect"
	"testing"

	backends "github.com/bennyz/example-finder/backend"
)

func TestDifference(t *testing.T) {
	arr1 := []int64{1, 2, 3, 4, 5}
	arr2 := []int64{4, 5, 7}
	got := Difference(arr1, arr2)
	expected := []int64{7}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Difference(arr1, arr2) = %v, want %v", got, expected)
	}
}

func TestSliceToString(t *testing.T) {
	arr := []int64{1, 2, 3}
	expected := "1,2,3"

	if got := SliceToString(arr); got != expected {
		t.Errorf("SliceToString(arr) = %v, want %v", got, expected)
	}
}

func TestMapToSlice(t *testing.T) {
	repo := backends.Result{
		RepoID:   1234,
		RepoName: "MyRepo",
		RepoURL:  "http://poop.com",
		Stars:    17,
	}

	m := map[int64]*backends.Result{1234: &repo}
	expected := []*backends.Result{&repo}

	if got := MapToSlice(m); !reflect.DeepEqual(got, expected) {
		t.Errorf("MapToSlice(m) = %v, want %v", got, expected)
	}
}
