package dbdiff

import (
	"reflect"
	"strings"
)

func AssertTypePtrOfSlice(ptr interface{}) bool {
	return reflect.ValueOf(ptr).Type().Elem().Kind() == reflect.Slice
}

// ptr the point of the input
func AssertTypePtrOfSliceWithStruct(ptr interface{}) bool {
	return AssertTypePtrOfSlice(ptr) &&
		reflect.ValueOf(ptr).Type().Elem().Elem().Kind() == reflect.Struct
}

func AssertTypePtrOfStruct(ptr interface{}) bool {
	return reflect.ValueOf(ptr).Type().Elem().Kind() == reflect.Struct
}

func AssertStrEmpty(str string) bool {
	return "" == str
}

func AssertStrBlank(str string) bool {
	return AssertStrEmpty(str) || len(strings.TrimSpace(str)) == 0
}
