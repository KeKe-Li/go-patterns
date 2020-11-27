package main

import "fmt"

type Country struct {
	Name string
}

type City struct {
	Name string
}
type Printable interface {
	PrintStr()
}

func (c Country) PrintStr() {
	fmt.Println(c.Name)
}

func (c City) PrintStr() {
	fmt.Println(c.Name)
}

func main() {
	c1 := Country{"keke"}
	c2 := City{"beijing"}

	c1.PrintStr()
	c2.PrintStr()
}
