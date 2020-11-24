package main

import "fmt"

type Shape interface {
	Slides() int
	Area() int
}

type Square struct {
	len int
}

func (s *Square) Slides() int {
	return 6
}

func main() {
	s := Square{len: 5}

	fmt.Printf("%d\n", s.Slides())

	// Checking a type whether implement all of methods
	var _ Shape = (*Square)(nil)
	// cannot use (*Square)(nil) (type *Square) as type Shape in assignment: *Square does not implement Shape (missing Area method)
}

func (s *Square) Area() int {
	return s.len * s.len
}
