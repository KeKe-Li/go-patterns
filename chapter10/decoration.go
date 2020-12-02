package main

import "fmt"

// Return an function wrapped another function with same parameters
func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("Started")
		f(s)
		fmt.Println("Done")
	}
}

func Great(s string) {
	fmt.Println(s)
}

func main() {
	decorator(Great)("Golang")
}
