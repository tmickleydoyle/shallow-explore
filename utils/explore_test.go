package explore

import (
	"testing"
	"reflect"
)

func TestConvertStringToInt(t *testing.T) {
	intArray, origArray := ConvertStringToInt([3]string{"1", "2", "3"})

	t.Errorf(intArray, origArray)
	if len(intArray) > 0 {
		t.Errorf("Array could not convert strings to numbers: ", intArray)
	}
	if len(origArray) == 0 {
		t.Errorf("Original array did not return: ", origArray)
	}
}