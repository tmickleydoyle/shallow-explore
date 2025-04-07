package explore

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

// CSVSummary generates a summary of the entire CSV file
func CSVSummary(records [][]string) map[string]interface{} {
	summary := make(map[string]interface{})
	
	// Basic file info
	summary["total_rows"] = len(records) - 1 // Excluding header
	summary["total_columns"] = len(records[0])
	summary["column_names"] = records[0]
	
	// Column type inference
	columnTypes := make(map[string]string)
	for i, colName := range records[0] {
		// Check first 10 rows (or fewer if less data) to infer type
		samples := []string{}
		for j := 1; j < min(len(records), 11); j++ {
			if i < len(records[j]) {
				samples = append(samples, records[j][i])
			}
		}
		columnTypes[colName] = inferColumnType(samples)
	}
	summary["column_types"] = columnTypes
	
	// Data completeness
	completeness := make(map[string]float64)
	for i, colName := range records[0] {
		filledCount := 0
		for j := 1; j < len(records); j++ {
			if i < len(records[j]) && records[j][i] != "" {
				filledCount++
			}
		}
		completeness[colName] = float64(filledCount) / float64(len(records)-1) * 100
	}
	summary["completeness"] = completeness
	
	return summary
}

// Helper function for CSVSummary
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Helper function to infer column type
func inferColumnType(samples []string) string {
	isInt := true
	isFloat := true
	
	for _, sample := range samples {
		if sample == "" {
			continue // Skip empty values for type inference
		}
		
		// Try parsing as int
		_, err := strconv.ParseInt(sample, 10, 64)
		if err != nil {
			isInt = false
		}
		
		// Try parsing as float
		_, err = strconv.ParseFloat(sample, 64)
		if err != nil {
			isFloat = false
		}
		
		if !isInt && !isFloat {
			break
		}
	}
	
	if isInt {
		return "integer"
	} else if isFloat {
		return "float"
	}
	return "string"
}

// FilterCSVData filters CSV data based on column and condition
func FilterCSVData(records [][]string, columnName string, condition string, value string) [][]string {
	filtered := [][]string{records[0]} // Start with header row
	
	// Find column index
	colIndex := -1
	for i, col := range records[0] {
		if col == columnName {
			colIndex = i
			break
		}
	}
	
	if colIndex == -1 {
		return filtered // Column not found
	}
	
	// Process each row
	for i := 1; i < len(records); i++ {
		row := records[i]
		if colIndex >= len(row) {
			continue // Skip rows with insufficient columns
		}
		
		cellValue := row[colIndex]
		includeRow := false
		
		switch condition {
		case "equals":
			includeRow = cellValue == value
		case "contains":
			includeRow = strings.Contains(cellValue, value)
		case "greater_than":
			cellFloat, err1 := strconv.ParseFloat(cellValue, 64)
			valueFloat, err2 := strconv.ParseFloat(value, 64)
			if err1 == nil && err2 == nil {
				includeRow = cellFloat > valueFloat
			}
		case "less_than":
			cellFloat, err1 := strconv.ParseFloat(cellValue, 64)
			valueFloat, err2 := strconv.ParseFloat(value, 64)
			if err1 == nil && err2 == nil {
				includeRow = cellFloat < valueFloat
			}
		case "starts_with":
			includeRow = strings.HasPrefix(cellValue, value)
		case "ends_with":
			includeRow = strings.HasSuffix(cellValue, value)
		default:
			includeRow = false
		}
		
		if includeRow {
			filtered = append(filtered, row)
		}
	}
	
	return filtered
}

// ExportToJSON exports CSV data to a JSON file
func ExportToJSON(records [][]string, outputPath string) error {
	result := []map[string]string{}
	headers := records[0]
	
	// Convert each row to a map
	for i := 1; i < len(records); i++ {
		row := make(map[string]string)
		for j, header := range headers {
			if j < len(records[i]) {
				row[header] = records[i][j]
			} else {
				row[header] = ""
			}
		}
		result = append(result, row)
	}
	
	// Marshal to JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	
	// Write to file
	return ioutil.WriteFile(outputPath, jsonData, 0644)
}

// CalculateCorrelation calculates Pearson correlation between two numeric columns
func CalculateCorrelation(records [][]string, column1 string, column2 string) (float64, error) {
	// Find column indices
	col1Index := -1
	col2Index := -1
	for i, col := range records[0] {
		if col == column1 {
			col1Index = i
		}
		if col == column2 {
			col2Index = i
		}
	}
	
	if col1Index == -1 || col2Index == -1 {
		return 0, fmt.Errorf("column not found")
	}
	
	// Extract numeric values
	var x []float64
	var y []float64
	
	for i := 1; i < len(records); i++ {
		row := records[i]
		if col1Index < len(row) && col2Index < len(row) {
			val1, err1 := strconv.ParseFloat(row[col1Index], 64)
			val2, err2 := strconv.ParseFloat(row[col2Index], 64)
			if err1 == nil && err2 == nil {
				x = append(x, val1)
				y = append(y, val2)
			}
		}
	}
	
	if len(x) < 2 {
		return 0, fmt.Errorf("insufficient numeric data for correlation")
	}
	
	// Calculate means
	xMean := 0.0
	yMean := 0.0
	for i := range x {
		xMean += x[i]
		yMean += y[i]
	}
	xMean /= float64(len(x))
	yMean /= float64(len(y))
	
	// Calculate correlation
	numerator := 0.0
	xDenom := 0.0
	yDenom := 0.0
	
	for i := range x {
		xDiff := x[i] - xMean
		yDiff := y[i] - yMean
		numerator += xDiff * yDiff
		xDenom += xDiff * xDiff
		yDenom += yDiff * yDiff
	}
	
	if xDenom == 0 || yDenom == 0 {
		return 0, nil // No variation in at least one variable
	}
	
	return numerator / math.Sqrt(xDenom * yDenom), nil
}

// DetectAnomalies detects anomalies in numeric columns using Z-score method
func DetectAnomalies(records [][]string, threshold float64) map[string][]int {
	anomolies := make(map[string][]int)
	
	// Process each column
	for colIndex, colName := range records[0] {
		// Extract numeric values
		values := []float64{}
		valueIndices := []int{}
		
		for i := 1; i < len(records); i++ {
			if colIndex < len(records[i]) {
				if val, err := strconv.ParseFloat(records[i][colIndex], 64); err == nil {
					values = append(values, val)
					valueIndices = append(valueIndices, i)
				}
			}
		}
		
		// Need enough data points for meaningful anomaly detection
		if len(values) < 5 {
			continue
		}
		
		// Calculate mean and standard deviation
		mean := 0.0
		for _, val := range values {
			mean += val
		}
		mean /= float64(len(values))
		
		stdDev := 0.0
		for _, val := range values {
			stdDev += math.Pow(val - mean, 2)
		}
		stdDev = math.Sqrt(stdDev / float64(len(values)))
		
		if stdDev == 0 {
			continue // Skip columns with no variation
		}
		
		// Find anomalies (values with Z-score above threshold)
		anomalousRows := []int{}
		for i, val := range values {
			zScore := math.Abs(val - mean) / stdDev
			if zScore > threshold {
				anomalousRows = append(anomalousRows, valueIndices[i])
			}
		}
		
		if len(anomalousRows) > 0 {
			anomolies[colName] = anomalousRows
		}
	}
	
	return anomolies
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