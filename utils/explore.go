package explore

import (
    "encoding/csv"
    "fmt"
    "log"
	"math"
    "os"
	"sort"
	"strconv"
	"strings"
)

type Sorted struct {
	Key   string
	Value int
}

// Create a type for a sorted list by values
type SortedList []Sorted

// This is used to find the len of the sorted list
func (p SortedList) Len() int           { return len(p) }
// This is used to reorder the elements in the list
func (p SortedList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
// This is used to compare the elements in the list
func (p SortedList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// Loads the CSV file
func ReadCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file "+filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for "+filePath, err)
    }

    return records
}

// Check the columns in the CSV file to converted string based int/float to float64
func ConvertStringToInt(stringArray []string) ([]float64, []string) {
    var intArray []float64

    for _, i := range stringArray {
        if convertNum, err := strconv.ParseFloat(i, 64); err == nil {
            intArray = append(intArray, convertNum)
        } else if intValue, err := strconv.ParseInt(i, 10, 64); err == nil {
			convertNum := float64(intValue)
            intArray = append(intArray, convertNum)
        }
    }

	return intArray, stringArray
}

// Finds the min and max from an array
func MinMaxValues(intArray []float64) (float64, float64) {
    var max float64 = intArray[0]
    var min float64 = intArray[0]
    for _, value := range intArray {

        if max < value {
            max = value
        }
        if min > value {
            min = value
        }
    }
    return min, max
}

// Find the mean value from an array
func MeanValue(intArray []float64) float64 {
	total := 0.0

	for _, v := range intArray {
		total += v
	}

	return math.Round(total / float64(len(intArray)))
}

// Find the median value from an array
func MedianValue(intArray []float64) float64 {
	sort.Float64s(intArray)
	mNumber := len(intArray) / 2

	if len(intArray)%2!=0 {
		return intArray[mNumber]
	}

	return (intArray[mNumber-1]+intArray[mNumber]) / 2
}

// Formatted string for number based columns
func FloatOutput(min float64, max float64, mean float64, median float64) string {
	outputString :=
`(The plot is a general trend of all points) 

   Min: %.2f
   Max: %.2f
  Mean: %.2f
Median: %.2f
`
	finalOutput := fmt.Sprintf(outputString, min, max, mean, median)
	return finalOutput
}

// Used to count how many times a string elements is in an array
func CountValues(stringArray []string) map[string]int {
	valuesMap := make(map[string]int)
	
	for _, v := range stringArray {
		v := string(v)
		if _, ok := valuesMap[v]; ok {
			valuesMap[v] = valuesMap[v] + 1
		} else {
			valuesMap[v] = 1
		}
	}

	return valuesMap
}

// Sorts a map by the value
func SortMapByValue(valueMap map[string]int) SortedList {
	sortedMap := make(SortedList, len(valueMap))

	i := 0
	for k, v := range valueMap {
		sortedMap[i] = Sorted{k, v}
		i++
	}

	sort.Sort(sortedMap)

	return sortedMap
}

// Builds a histogram and final output for string based columns
func HistTopTen(sortedList SortedList, column string) string {
	var key string
	var barValue int
	max := float64(sortedList[0].Value)
	histString := "Horizontal Histogram"

	if len(sortedList) > 10 {
		histString = histString + " - Top Ten\n\n"
	} else {
		histString = histString + "\n\n"
	}

	i := 0
	for d := range sortedList {
		if i < 10 {
			if max > 75 {
				barValue = int((float64(sortedList[d].Value) / max) * 75)
			} else {
				barValue = sortedList[d].Value
			}
			bar := strings.Repeat("*", barValue)
			if len(sortedList[d].Key) > 20 {
				key = sortedList[d].Key[:17]+"..."
			} else {
				key = sortedList[d].Key
			}
			histString = histString + fmt.Sprintf("%20v: %s (%d)\n", key, bar, sortedList[d].Value)
		}
		i++

	}

	histString = histString + fmt.Sprintf("%20v: %d\n", "Unique Strings", len(sortedList))

	return histString
}
