package dbdiff

import (
	"testing"
)

type A struct {
}

func TestAssertTypePtrOfSliceWithStruct(t *testing.T) {
	in1 := []A{}
	verify(t, 1, "AssertTypePtrOfSliceWithStruct", &in1, AssertTypePtrOfSliceWithStruct(&in1), true)
	in2 := 2
	verify(t, 2, "AssertTypePtrOfSliceWithStruct", &in2, AssertTypePtrOfSliceWithStruct(&in2), false)
}

func TestAssertTypePtrOfStruct(t *testing.T) {
	a := A{}
	verify(t, 1, "AssertTypePtrOfStruct", &a, AssertTypePtrOfStruct(&a), true)
}

func verify(t *testing.T, testnum int, testcase string, input, output, expected interface{}) {
	if expected != output {
		t.Errorf("%d. %s with input = %v: output %v != %v", testnum, testcase, input, output, expected)
	}
}
