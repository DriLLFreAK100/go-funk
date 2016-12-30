package fn

import (
	"fmt"
	"reflect"
	"strings"
)

func equal(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)

}

// Contains is ...
func Contains(in interface{}, elem interface{}) bool {
	inValue := reflect.ValueOf(in)

	elemValue := reflect.ValueOf(elem)

	inType := inValue.Type()

	if inType.Kind() == reflect.String {
		return strings.Contains(inValue.String(), elemValue.String())
	}

	if inType.Kind() == reflect.Map {
		keys := inValue.MapKeys()
		for i := 0; i < len(keys); i++ {
			if equal(keys[i].Interface(), elem) {
				return true
			}
		}
	}

	if inType.Kind() == reflect.Slice {
		for i := 0; i < inValue.Len(); i++ {
			if equal(inValue.Index(i).Interface(), elem) {
				return true
			}
		}

	}

	return false
}

// ToMap transforms a slice of instances to a Map
// []*Foo => Map<int, *Foo>
func ToMap(in interface{}, pivot string) interface{} {
	value := reflect.ValueOf(in)

	// input value must be a slice
	if value.Kind() != reflect.Slice {
		panic(fmt.Sprintf("%v must be a slice", in))
	}

	inType := value.Type()

	structType := inType.Elem()

	// retrieve the struct in the slice to deduce key type
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}

	field, _ := structType.FieldByName(pivot)

	// value of the map will be the input type
	collectionType := reflect.MapOf(field.Type, inType.Elem())

	// create a map from scratch
	collection := reflect.MakeMap(collectionType)

	for i := 0; i < value.Len(); i++ {
		instance := value.Index(i)
		var field reflect.Value

		if instance.Kind() == reflect.Ptr {
			field = instance.Elem().FieldByName(pivot)
		} else {
			field = instance.FieldByName(pivot)
		}

		collection.SetMapIndex(field, instance)
	}

	return collection.Interface()
}

// Map is ...
func Map(arr interface{}, mapFunc interface{}) interface{} {
	funcValue := reflect.ValueOf(mapFunc)
	arrValue := reflect.ValueOf(arr)

	// Retrieve the type, and check if it is one of the array or slice.
	arrType := arrValue.Type()
	arrElemType := arrType.Elem()
	if arrType.Kind() != reflect.Array && arrType.Kind() != reflect.Slice {
		panic("Array parameter's type is neither array nor slice.")
	}

	funcType := funcValue.Type()

	// Checking whether the second argument is function or not.
	// And also checking whether its signature is func ({type A}) {type B}.
	if funcType.Kind() != reflect.Func || funcType.NumIn() != 1 || funcType.NumOut() != 1 {
		panic("Second argument must be map function.")
	}

	// Checking whether element type is convertible to function's first argument's type.
	if !arrElemType.ConvertibleTo(funcType.In(0)) {
		panic("Map function's argument is not compatible with type of array.")
	}

	// Get slice type corresponding to function's return value's type.
	resultSliceType := reflect.SliceOf(funcType.Out(0))

	// MakeSlice takes a slice kind type, and makes a slice.
	resultSlice := reflect.MakeSlice(resultSliceType, 0, arrValue.Len())

	for i := 0; i < arrValue.Len(); i++ {
		resultSlice = reflect.Append(resultSlice, funcValue.Call([]reflect.Value{arrValue.Index(i)})[0])
	}

	// Convering resulting slice back to generic interface.
	return resultSlice.Interface()
}

// MakeSlice is ...
func MakeSlice(in interface{}) interface{} {
	value := reflect.ValueOf(in)

	sliceType := reflect.SliceOf(reflect.TypeOf(in))
	slice := reflect.New(sliceType)
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)
	sliceValue = reflect.Append(sliceValue, value)
	slice.Elem().Set(sliceValue)

	return slice.Elem().Interface()
}

// MakeZero is ...
func MakeZero(in interface{}) interface{} {
	return reflect.Zero(reflect.TypeOf(in)).Interface()
}