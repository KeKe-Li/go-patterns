package chapter11

import "fmt"

type DecoratedVisitor struct {
	visitor    Visitor
	decorators []VisitorFunc
}

func NewDecoratedVisitor(v Visitor, fn ...VisitorFunc) Visitor {
	if len(fn) == 0 {
		return v
	}
	return DecoratedVisitor{v, fn}
}

// Visit implements Visitor
func (v DecoratedVisitor) Visit(fn VisitorFunc) error {
	return v.visitor.Visit(func(info *Info, err error) error {
		if err != nil {
			return err
		}
		if err := fn(info, nil); err != nil {
			return err
		}
		for i := range v.decorators {
			if err := v.decorators[i](info, nil); err != nil {
				return err
			}
		}
		return nil
	})
}

func NameVisitors(info *Info, err error) error {
	fmt.Printf("Name=%s,Namespace=%s\n", info.Name, info.Namespace)
	return nil
}

func OtherVisitor(info *Info, err error) error {
	fmt.Printf("Other=%s\n", info.OtherThings)
	return nil
}

func LoadFile(info *Info, err error) error {
	info.Name = "keke"
	info.Namespace = "jame"
	info.OtherThings = " golang language"
	return nil
}

func main() {
	info := Info{}
	var v Visitor = &info
	v = NewDecoratedVisitor(v, NameVisitors, OtherVisitor)
	v.Visit(LoadFile)
}
