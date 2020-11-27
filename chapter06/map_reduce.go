package main

import (
	"fmt"
	"strings"
)

func MapUpcase(arr []string, fn func(s string) string) []string {
	var newArray = []string{}
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func MapLen(arr []string, fn func(s string) int) []int {
	var newArray = []int{}
	for _, it := range arr {
		newArray = append(newArray, fn(it))
	}
	return newArray
}

func Reduce(arr []string, fn func(s string) int) int {
	sum := 0
	for _, it := range arr {
		sum += fn(it)
	}
	return sum
}

func Filter(arr []int, fn func(n int) bool) []int {
	var newArray = []int{}
	for _, it := range arr {
		if fn(it) {
			newArray = append(newArray, it)
		}
	}
	return newArray
}

func main() {
	var list = []string{"keke", "li", "ruan", "lv"}
	val := MapUpcase(list, func(s string) string {
		return strings.ToUpper(s)
	})
	fmt.Printf("%v\n", val)

	var list2 = []string{"Golang", "Great"}
	val2 := Reduce(list2, func(s string) int {
		return len(s)
	})

	fmt.Printf("%v\n", val2)

	var inset = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	out := Filter(inset, func(n int) bool {
		return n%2 == 1
	})

	fmt.Printf("%v\n", out)
}
