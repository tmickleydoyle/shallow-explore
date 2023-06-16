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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Sorted struct {
	Key   string
	Value int
}

type SortedList []Sorted

func (p SortedList) Len() int           { return len(p) }
func (p SortedList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p SortedList) Less(i, j int) bool { return p[i].Value > p[j].Value }

func ReadCSVFile(filePath string) [][]string {
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

func MeanValue(intArray []float64) float64 {
	total := 0.0

	for _, v := range intArray {
		total += v
	}

	return math.Round(total / float64(len(intArray)))
}

func MedianValue(intArray []float64) float64 {
	sort.Float64s(intArray)
	mNumber := len(intArray) / 2

	if len(intArray)%2 != 0 {
		return intArray[mNumber]
	}

	return (intArray[mNumber-1] + intArray[mNumber]) / 2
}

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
			bar := strings.Repeat("☐", barValue)
			if len(sortedList[d].Key) > 20 {
				key = sortedList[d].Key[:17] + "..."
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

type ChatGPTRequest struct {
	Message string `json:"message"`
}

type ChatGPTResponse struct {
	Reply string `json:"reply"`
}

func CallChatGPTAPI(message string, authToken string) (string, error) {
	payload := ChatGPTRequest{
		Message: message,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.chatgpt.com/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var apiResponse ChatGPTResponse
	err = json.Unmarshal(respBody, &apiResponse)
	if err != nil {
		return "", err
	}

	return apiResponse.Reply, nil
}