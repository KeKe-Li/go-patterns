package chapter09

import "reflect"

// Reduce computes the reduction of the pair function across the elements of
// the slice.
func Reduce(slice, pairFunction, zero interface{}) interface{} {
	in := reflect.ValueOf(slice)
	if in.Kind() != reflect.Slice {
		panic("reduce: not slice")
	}

	n := in.Len()
	switch n {
	case 0:
		return zero
	case 1:
		return in.Index(0)
	}
	elemType := in.Type().Elem()
	fn := reflect.ValueOf(pairFunction)
	if !goodFunc(fn, elemType, elemType, elemType) {
		str := elemType.String()
		panic("apply: function must be of type func(" + str + ", " + str + ") " + str)
	}
	// DO the first two by hand to prime the pump
	var ins [2]reflect.Value
	ins[0] = in.Index(0)
	ins[1] = in.Index(1)
	out := fn.Call(ins[:])[0]
	for i := 2; i < n; i++ {
		ins[0] = out
		ins[1] = in.Index(i)
		out = fn.Call(ins[:])[0]
	}
	return out.Interface()
}
