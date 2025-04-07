package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
	"github.com/charmbracelet/lipgloss"
	explore "github.com/tmickleydoyle/shallow-explore/utils"	
)

var (
	file      string
	path      string
	style     string
	filter    string
	column    string
	condition string
	value     string
	exportJson bool
	exportPath string
	correlate  bool
	col1       string
	col2       string
	summary    bool
	anomalies  bool
	threshold  float64
)

var styleDark = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#808080")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(2).
	PaddingRight(2).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#FAFAFA"))

var styleLight = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#808080")).
	Background(lipgloss.Color("#FAFAFA")).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(2).
	PaddingRight(2).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#808080"))

func main() {
	// Basic flags
	flag.StringVar(&file, "csv", "", "path to CSV file")
	flag.StringVar(&style, "style", "", "output style (dark or light)")
	
	// Data filtering flags
	flag.StringVar(&filter, "filter", "", "enable filtering mode")
	flag.StringVar(&column, "column", "", "column to filter on")
	flag.StringVar(&condition, "condition", "equals", "filter condition (equals, contains, greater_than, less_than, starts_with, ends_with)")
	flag.StringVar(&value, "value", "", "filter value to compare against")
	
	// JSON export flags
	flag.BoolVar(&exportJson, "export-json", false, "export data to JSON")
	flag.StringVar(&exportPath, "export-path", "", "path for exported JSON file")
	
	// Correlation flags
	flag.BoolVar(&correlate, "correlate", false, "calculate correlation between two columns")
	flag.StringVar(&col1, "col1", "", "first column for correlation")
	flag.StringVar(&col2, "col2", "", "second column for correlation")
	
	// Summary flag
	flag.BoolVar(&summary, "summary", false, "generate summary of CSV file")
	
	// Anomaly detection
	flag.BoolVar(&anomalies, "anomalies", false, "detect anomalies in numeric columns")
	flag.Float64Var(&threshold, "threshold", 3.0, "z-score threshold for anomaly detection (default: 3.0)")
	
	flag.Parse()

	if file == "" {
		log.Fatal("Could not find the path to the CSV file")
	}
	path = file
	records := explore.ReadCSVFile(path)
	
	// Set default style
	var selectedStyle lipgloss.Style
	if style == "dark" {
		selectedStyle = styleDark
	} else {
		selectedStyle = styleLight
	}
	
	// Handle summary mode
	if summary {
		summaryData := explore.CSVSummary(records)
		fmt.Println(selectedStyle.Render(fmt.Sprintf("CSV Summary for: %s\n", path)))
		
		// Display summary information
		fmt.Println(selectedStyle.Render(fmt.Sprintf("Total Rows: %d\n", summaryData["total_rows"])))
		fmt.Println(selectedStyle.Render(fmt.Sprintf("Total Columns: %d\n", summaryData["total_columns"])))
		
		// Display column names and types
		fmt.Println(selectedStyle.Render("Column Information:"))
		columnTypes := summaryData["column_types"].(map[string]string)
		completeness := summaryData["completeness"].(map[string]float64)
		
		for _, colName := range summaryData["column_names"].([]string) {
			fmt.Println(selectedStyle.Render(fmt.Sprintf("  %s (Type: %s, Completeness: %.1f%%)", 
				colName, columnTypes[colName], completeness[colName])))
		}
		return
	}
	
	// Handle filtering
	if filter != "" && column != "" && value != "" {
		records = explore.FilterCSVData(records, column, condition, value)
		fmt.Println(selectedStyle.Render(fmt.Sprintf("Filtered data where %s %s %s (Found %d records)\n", 
			column, condition, value, len(records)-1)))
	}
	
	// Handle JSON export
	if exportJson {
		if exportPath == "" {
			// Generate default filename based on input file if not provided
			base := filepath.Base(path)
			baseName := strings.TrimSuffix(base, filepath.Ext(base))
			exportPath = baseName + "_" + time.Now().Format("20060102_150405") + ".json"
		}
		
		err := explore.ExportToJSON(records, exportPath)
		if err != nil {
			log.Fatalf("Failed to export to JSON: %v", err)
		}
		
		fmt.Println(selectedStyle.Render(fmt.Sprintf("Data exported to JSON: %s\n", exportPath)))
		return
	}
	
	// Handle correlation calculation
	if correlate && col1 != "" && col2 != "" {
		correlation, err := explore.CalculateCorrelation(records, col1, col2)
		if err != nil {
			log.Fatalf("Failed to calculate correlation: %v", err)
		}
		
		fmt.Println(selectedStyle.Render(fmt.Sprintf("Correlation between '%s' and '%s': %.4f\n", 
			col1, col2, correlation)))
		
		// Interpretation guide
		var interpretation string
		corrAbs := math.Abs(correlation)
		
		if corrAbs >= 0.9 {
			interpretation = "Very strong relationship"
		} else if corrAbs >= 0.7 {
			interpretation = "Strong relationship"
		} else if corrAbs >= 0.5 {
			interpretation = "Moderate relationship"
		} else if corrAbs >= 0.3 {
			interpretation = "Weak relationship"
		} else {
			interpretation = "Little to no relationship"
		}
		
		if correlation < 0 {
			interpretation += " (negative/inverse)"
		} else {
			interpretation += " (positive/direct)"
		}
		
		fmt.Println(selectedStyle.Render("Interpretation: " + interpretation))
		return
	}
	
	// Handle anomaly detection
	if anomalies {
		anomalousRows := explore.DetectAnomalies(records, threshold)
		
		if len(anomalousRows) == 0 {
			fmt.Println(selectedStyle.Render(fmt.Sprintf("No anomalies detected with threshold %.1f\n", threshold)))
		} else {
			fmt.Println(selectedStyle.Render(fmt.Sprintf("Anomalies detected with threshold %.1f:\n", threshold)))
			
			for colName, rows := range anomalousRows {
				fmt.Println(selectedStyle.Render(fmt.Sprintf("Column '%s': %d anomalies found", colName, len(rows))))
				
				// Show first 5 anomalies at most
				displayRows := rows
				if len(rows) > 5 {
					displayRows = rows[:5]
				}
				
				// Get column index
				colIndex := -1
				for i, col := range records[0] {
					if col == colName {
						colIndex = i
						break
					}
				}
				
				for _, rowIdx := range displayRows {
					if colIndex < len(records[rowIdx]) {
						fmt.Println(selectedStyle.Render(fmt.Sprintf("  Row %d: %s", rowIdx, records[rowIdx][colIndex])))
					}
				}
				
				if len(rows) > 5 {
					fmt.Println(selectedStyle.Render(fmt.Sprintf("  ... and %d more anomalies", len(rows)-5)))
				}
			}
		}
		return
	}
	
	// Default behavior: explore each column
	for column := range records[0] {
		colValues := []string{}

		for i := 1; i < len(records); i++ {
			colValues = append(colValues, records[i][column])
		}

		transformedArray, _ := explore.ConvertStringToInt(colValues)
		plotArray, stringValues := explore.ConvertStringToInt(colValues)
		column := fmt.Sprintf("Column: %s\n\n", records[0][column])

		if len(transformedArray) > 0 {
			min, max := explore.MinMaxValues(transformedArray)
			mean := explore.MeanValue(transformedArray)
			median := explore.MedianValue(transformedArray)
			statsOutput := explore.FloatOutput(min, max, mean, median)
			graph := asciigraph.Plot(plotArray, asciigraph.Height(20), asciigraph.Width(90), asciigraph.Caption(statsOutput))
			fmt.Println(selectedStyle.Render(column + graph))
		} else {
			valuesMap := explore.CountValues(stringValues)
			sortedMap := explore.SortMapByValue(valuesMap)
			histogram := explore.HistTopTen(sortedMap, column)
			fmt.Println(selectedStyle.Render(column + histogram))
		}
	}
}
