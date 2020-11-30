package chater07

import "reflect"

/*
• Sensible defaults

• Highly configurable

• Easy to maintain

• Self documenting

• Safe for newcomers

• No need nil or an empty value
*/

func Transform(slice, fn interface{}) interface{} {
	return transform(slice, fn, false)
}

func TransformInPlace(slice, fn interface{}) interface{} {
	return transform(slice, fn, true)
}

func transform(slice, function interface{}, inPlace bool) interface{} {
	// check the `slice` type is Slice
	sliceInType := reflect.ValueOf(slice)
	if sliceInType.Kind() != reflect.Slice {
		panic("transform: not slice")
	}
	// check the function signature
	fn := reflect.ValueOf(function)
	elemType := sliceInType.Type().Elem()
	if !verifyFuncSignature(fn, elemType, nil) {
		panic("transform: function must be of type func(" + sliceInType.Type().Elem().String() + ") outputElemType")
	}
	sliceOutType := sliceInType
	if !inPlace {
		sliceOutType = reflect.MakeSlice(reflect.SliceOf(fn.Type().Out(0)), sliceInType.Len(), sliceInType.Len())
	}
	for i := 0; i < sliceInType.Len(); i++ {
		sliceOutType.Index(i).Set(fn.Call([]reflect.Value{sliceInType.Index(i)})[0])
	}
	return sliceOutType.Interface()
}

// Verify Function Signature
func verifyFuncSignature(fn reflect.Value, types ...reflect.Type) bool {
	// check it is a function
	if fn.Kind() != reflect.Func {
		return false
	}
	// NumIn() - returns a function type's input parameter count.
	// NumOut() - returns a function type's output parameter count.
	if (fn.Type().NumIn() != len(types)-1) || (fn.Type().NumOut() != 1) {
		return false
	}
	// In() - returns the type of a function types's i'th input parameter
	for i := 0; i < len(types)-1; i++ {
		if fn.Type().In(i) != types[i] {
			return false
		}
	}
	// Out() - returns the type of a function type's i'th output parameter
	outType := types[len(types)-1]
	if outType != nil && fn.Type().Out(0) != outType {
		return false
	}
	return true
}
