package main

import (
	"flag"
    "fmt"
	"log"

	"github.com/guptarohit/asciigraph"
	"github.com/charmbracelet/lipgloss"
	explore "github.com/tmickleydoyle/shallow-explore/utils"
)

var (
	file string
	path string
)

// This defines the look of the printed content to the terminal.
var style = lipgloss.NewStyle().
    Bold(true).
    Foreground(lipgloss.Color("#FAFAFA")).
    Background(lipgloss.Color("#808080")).
    PaddingTop(1).
	PaddingBottom(1).
    PaddingLeft(2).
	PaddingRight(2).
	BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#FAFAFA"))

// The final logic for pull together outputs for each column
func main() {
	flag.StringVar(&file, "csv", "", "starting point")
	flag.Parse()

	if file != "" {
		path = file
	} else {
		log.Fatal("Could not find the path to the CSV file")
	}

    records := explore.ReadCsvFile(path)
	for column := range records[0] {
		colValues := []string{}

		for i := 1; i < len(records); i++ {
			colValues = append(colValues, records[i][column])
		}

		transformedArray, _ := explore.ConvertStringToInt(colValues)
		plotArray, stringValues := explore.ConvertStringToInt(colValues)
		column := fmt.Sprintf("Column: %s\n\n", records[0][column])
		
		// The first part of the if is for int/float columns, and the
		// else is used for string based columns
		if len(transformedArray) > 0 {
			min, max := explore.MinMaxValues(transformedArray)
			mean := explore.MeanValue(transformedArray)
			median := explore.MedianValue(transformedArray)
			statsOutput := explore.FloatOutput(min, max, mean, median)
			graph := asciigraph.Plot(plotArray, asciigraph.Height(20), asciigraph.Width(90), asciigraph.Caption(statsOutput))
			fmt.Println(style.Render(column + graph))
		} else {
			valuesMap := explore.CountValues(stringValues)
			sortedMap := explore.SortMapByValue(valuesMap)
			histogram := explore.HistTopTen(sortedMap, column)
			fmt.Println(style.Render(column + histogram))
		}
	}
}