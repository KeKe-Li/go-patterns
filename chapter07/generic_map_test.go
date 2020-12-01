package chater07

import (
	"reflect"
	"testing"
)

func TestTransformString(t *testing.T) {
	list := []string{"1", "2", "3", "4", "5", "6"}
	expect := []string{"111", "222", "333", "444", "555", "666"}
	result := Transform(list, func(a string) string {
		return a + a + a
	})
	if !reflect.DeepEqual(expect, result) {
		t.Fatalf("Transform failed: expect %v got %v", expect, result)
	}
}

func TestTransformnPlace(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect := []int{3, 6, 9, 12, 15, 18, 21, 24, 27}
	TransformInPlace(list, func(a int) int {
		return a * 3
	})
	if !reflect.DeepEqual(expect, list) {
		t.Fatalf("Transform failed: expect %v got %v", expect, list)
	}
}

func TestTransformInPlace(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expect := []int{3, 6, 9, 12, 15, 18, 21, 24, 27}
	TransformInPlace(list, func(a int) int {
		return a * 3
	})
	if reflect.DeepEqual(expect, list) {
		t.Fatalf("Transform failed: expect %v got %v", expect, list)
	}
}
