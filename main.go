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
	file  string
	path  string
	style string
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
	flag.StringVar(&file, "csv", "", "starting point")
	flag.StringVar(&style, "style", "", "output style (dark or light)")
	flag.Parse()

	message := "Hello, ChatGPT!"

	reply, err := CallChatGPTAPI(message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)

	if file != "" {
		path = file
	} else {
		log.Fatal("Could not find the path to the CSV file")
	}

	records := ReadCSVFile(path)
	for column := range records[0] {
		colValues := []string{}

		for i := 1; i < len(records); i++ {
			colValues = append(colValues, records[i][column])
		}

		transformedArray, _ := ConvertStringToInt(colValues)
		plotArray, stringValues := ConvertStringToInt(colValues)
		column := fmt.Sprintf("Column: %s\n\n", records[0][column])

		var selectedStyle lipgloss.Style
		if style == "dark" {
			selectedStyle = styleDark
		} else {
			selectedStyle = styleLight
		}

		if len(transformedArray) > 0 {
			min, max := MinMaxValues(transformedArray)
			mean := MeanValue(transformedArray)
			median := MedianValue(transformedArray)
			statsOutput := FloatOutput(min, max, mean, median)
			graph := asciigraph.Plot(plotArray, asciigraph.Height(20), asciigraph.Width(90), asciigraph.Caption(statsOutput))
			fmt.Println(selectedStyle.Render(column + graph))
		} else {
			valuesMap := CountValues(stringValues)
			sortedMap := SortMapByValue(valuesMap)
			histogram := HistTopTen(sortedMap, column)
			fmt.Println(selectedStyle.Render(column + histogram))
		}
	}
}
