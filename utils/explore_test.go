package explore

import (
	"testing"
)

func TestConvertStringToInt(t *testing.T) {
	intArray, origArray := ConvertStringToInt([]string{"1", "2", "3"})

	if len(intArray) == 0 {
		t.Errorf("Array could not convert strings to numbers")
	}
	if len(origArray) == 0 {
		t.Errorf("Original array did not return: ")
	}

	stringArray, _ := ConvertStringToInt([]string{"one", "two", "three"})

	if len(stringArray) != 0 {
		t.Errorf("Incorrectly converted string to float")
	}
}

func TestMinMaxValues(t *testing.T) {
	minValue, maxValue := MinMaxValues([]float64{1.0, 2.0, 3.0})

	if minValue != 1 {
		t.Errorf("Min value is not calculated properly")
	}
	if maxValue != 3 {
		t.Errorf("Min value is not calculated properly")
	}
}

func TestMeanValue(t *testing.T) {
	meanValue := MeanValue([]float64{1.0, 2.0, 3.0})

	if meanValue != 2 {
		t.Errorf("Mean value is not calculated properly")
	}
}

func TestMedianValue(t *testing.T) {
	medianValue := MedianValue([]float64{1.0, 2.0, 3.0})

	if medianValue != 2 {
		t.Errorf("Mean value is not calculated properly")
	}
}
