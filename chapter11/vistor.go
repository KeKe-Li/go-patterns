package chapter11

import (
	"fmt"
)

type Info struct {
	Namespace   string
	Name        string
	OtherThings string
}

type VisitorFunc func(*Info, error) error

type Visitor interface {
	Visit(VisitorFunc) error
}

func (info *Info) Visit(fn VisitorFunc) error {
	return fn(info, nil)
}

type OtherThingsVisitor struct {
	visitor Visitor
}

func (o OtherThingsVisitor) Visit(fn VisitorFunc) error {
	return o.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("OtherThingsVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("===> OtherThings=%s\n", info.OtherThings)
		}
		fmt.Println("OtherThingsVisitor() after call function ")
		return err
	})
}

type LogVisitor struct {
	visitor Visitor
}

func (o LogVisitor) Visit(fn VisitorFunc) error {
	return o.visitor.Visit(func(info *Info, err error) error {
		fmt.Println("LogVisitor() before call function")
		err = fn(info, err)
		fmt.Println("LogVisitor() after call function")
		return err
	})
}

type NameVisitor struct {
	vistor Visitor
}

func (v NameVisitor) Visit(fn VisitorFunc) error {
	return v.vistor.Visit(func(info *Info, err error) error {
		fmt.Println("NameVisitor() before call function")
		err = fn(info, err)
		if err == nil {
			fmt.Printf("===> Name=%s,NameSpace=%s\n", info.Name, info.Namespace)
		}
		fmt.Println("NameVisitor() after call function")
		return err
	})
}
