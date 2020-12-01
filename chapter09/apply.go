package chapter09

import "reflect"

// Apply takes a slice of type []T and a function of type func(T) T. (If the
// input conditions are not satisfied, Apply panics.) It returns a newly
// allocated slice where each element is the result of calling the function on
// successive elements of the slice.
func Apply(slice, function interface{}) interface{} {
	return apply(slice, function, false)
}

// ApplyInPlace is like Apply, but overwrites the slice rather than returning a
// newly allocated slice.
func ApplyInPlace(slice, function interface{}) {
	apply(slice, function, true)
}

// Choose takes a slice of type []T and a function of type func(T) bool. (If
//// the input conditions are not satisfied, Choose panics.) It returns a newly
//// allocated slice containing only those elements of the input slice that
//// satisfy the function.
func Choose(slice, function interface{}) interface{} {
	out, _ := chooseOrDrop(slice, function, false, true)
	return out
}

func Drop(slice, function interface{}) interface{} {
	out, _ := chooseOrDrop(slice, function, false, false)
	return out
}

func ChooseInPalce(pointerToSlice, function interface{}) {
	chooseOrDropInplace(pointerToSlice, function, true)
}

func DropInPlace(pointerToSlice, function interface{}) {
	chooseOrDropInplace(pointerToSlice, function, false)
}

func apply(slice, function interface{}, inPalce bool) interface{} {
	// special case for strings,very common.
	if strSlice, ok := slice.([]string); ok {
		if strFn, ok := function.(func(string) string); ok {
			r := strSlice
			if !inPalce {
				r = make([]string, len(strSlice))
			}
			for i, s := range strSlice {
				r[i] = strFn(s)
			}
		}
	}
	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("apply: not slice")
	}
	fn := reflect.ValueOf(function)
	elemType := in.Type().Elem()
	if !goodFunc(fn, elemType, nil) {
		panic("apply:function must be type func(" + in.Type().Elem().String() + ") outputElemType")
	}
	out := in
	if !inPalce {
		out = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), in.Len(), in.Cap())
	}
	var ins [1]reflect.Value //Outside the loop to aovid one allocation
	for i := 0; i < in.Len(); i++ {
		ins[0] = in.Index(i)
		out.Index(i).Set(fn.Call(ins[:])[0])
	}
	return out.Interface()
}

func chooseOrDropInplace(slice, function interface{}, truth bool) {
	inp := reflect.ValueOf(slice)
	if inp.Kind() != reflect.Ptr {
		panic("choose or drop: not pointer to slice")
	}
	_, n := chooseOrDrop(inp.Elem().Interface(), function, true, truth)
	inp.Elem().SetLen(n)
}

var boolType = reflect.ValueOf(true).Type()

func chooseOrDrop(slice, function interface{}, inPlace, truth bool) (interface{}, int) {
	// Special case for strings, very common
	if strSlice, ok := slice.([]string); ok {
		if strFn, ok := function.(func(string) bool); ok {
			var r []string
			if inPlace {
				r = strSlice[:0]
			}
			for _, s := range strSlice {
				if strFn(s) == truth {
					r = append(r, s)
				}
			}
			return r, len(r)
		}
	}
	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("choose or Drop: not slice")
	}
	fn := reflect.ValueOf(function)
	elemType := in.Type().Elem()
	if !goodFunc(fn, elemType, boolType) {
		panic("choose/drop: function must be of type func(" + elemType.String() + ") bool")
	}
	var which []int
	var ins [1]reflect.Value //Outside the loop to avoid one allocation.
	for i := 0; i < in.Len(); i++ {
		ins[0] = in.Index(i)
		if fn.Call(ins[:])[0].Bool() == truth {
			which = append(which, i)
		}
	}
	out := in
	if !inPlace {
		out = reflect.MakeSlice(in.Type(), len(which), len(which))
	}
	for i := range which {
		out.Index(i).Set(in.Index(which[i]))
	}
	return out.Interface(), len(which)
}

// goodFunc verifies that the function satisfies the signature, represented as a slice of types.
func goodFunc(fn reflect.Value, types ...reflect.Type) bool {
	if fn.Kind() != reflect.Func {
		return false
	}
	// last type is return ,the rest are ins.
	if fn.Type().NumIn() != len(types)-1 || fn.Type().NumOut() != 1 {
		return false
	}
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}
