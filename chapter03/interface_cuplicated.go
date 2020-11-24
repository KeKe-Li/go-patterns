package main

import "fmt"

// Using a shared struct
type WithName struct {
	Name string
}

type Province struct {
	WithName
}

type Name struct {
	WithName
}

type Printables interface {
	PrintStr()
}

func (w WithName) Printables(){
	fmt.Println(w.Name)
}

func main(){
	c1 := Province{WithName{Name:"henan"}}
	c2 := Name{WithName{Name:"keke"}}

	c1.Printables()
	c2.Printables()
}
