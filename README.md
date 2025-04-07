# shallow-explore
From the command line, quickly explore data from a CSV file.

`shallow-explore` is a [Golang](https://go.dev/) backed command-line tool for iterating over columns from a CSV file. This is a gut check tool to make sure the assumptions about the data are within the expected range of normal.

## How-To

After installation, run the following command to start analyzing data:

```bash
# Style (default): light mode
shallow-explore -csv ~/complete/path/to/file/sample.csv

# Style: dark mode
shallow-explore -csv ~/complete/path/to/file/sample.csv -style dark

# Style: light mode
shallow-explore -csv ~/complete/path/to/file/sample.csv -style light
```

Note: The complete path of the file is required to load the data into the program.

### New Features

#### Data Summary
Generate a comprehensive summary of your CSV file with statistics about each column:

```bash
shallow-explore -csv ~/path/to/file/sample.csv -summary
```

#### Data Filtering
Filter data based on column values:

```bash
# Filter rows where the "Age" column equals "30"
shallow-explore -csv ~/path/to/file/sample.csv -filter true -column "Age" -condition equals -value "30"

# Available conditions: equals, contains, greater_than, less_than, starts_with, ends_with
shallow-explore -csv ~/path/to/file/sample.csv -filter true -column "Name" -condition contains -value "Smith"
```

#### Export to JSON
Export your CSV data to JSON format:

```bash
shallow-explore -csv ~/path/to/file/sample.csv -export-json -export-path "output.json"

# The export-path is optional. If not provided, a timestamped filename will be generated
shallow-explore -csv ~/path/to/file/sample.csv -export-json
```

#### Data Correlation
Calculate the correlation between two numeric columns:

```bash
shallow-explore -csv ~/path/to/file/sample.csv -correlate -col1 "Height" -col2 "Weight"
```

#### Anomaly Detection
Detect anomalies in numeric columns using Z-score method:

```bash
# Default threshold is 3.0
shallow-explore -csv ~/path/to/file/sample.csv -anomalies

# Custom threshold
shallow-explore -csv ~/path/to/file/sample.csv -anomalies -threshold 2.5
```

### Output

`shallow-explore` supports three types of data: integers, floats, and strings.

The following output is an example of an integer or float column. The column name at the top of the frame followed by a summary line graph of the items, and some quick statistics about the data.

<img width="709" alt="Screen Shot 2022-01-11 at 8 31 11 PM" src="https://user-images.githubusercontent.com/8069675/149228948-2dc71027-858e-406c-b09b-65231c9c04ca.png">

For string-based data, the column name is still at the top of the output. Below the column name lives a horizontal histogram and a count of unique entities found in the column.

<img width="769" alt="Screen Shot 2022-01-11 at 8 30 47 PM" src="https://user-images.githubusercontent.com/8069675/149228970-7cebd181-4faa-4369-886d-8e58650fca81.png">

## Installation

If Golang is installed, run the following command:

```bash
go install github.com/tmickleydoyle/shallow-explore
```

## Instructions for Installing Go

[Go docs](https://go.dev/)

### Installation with Homebrew

```bash
brew install go
```

## Why I Built This Tool

I find myself running and rerunning the same basic statistical analysis on data to get an understanding of how trends are moving. I figured why not make it easier and share it with everyone else! I hope this speeds up your decision making :heart:

## Feature Overview

- **CSV Data Exploration**: Visualize and analyze CSV data with automatic recognition of data types
- **Data Summary**: Get a comprehensive overview of your CSV data with column types and completeness percentages
- **Data Filtering**: Filter CSV data based on various conditions to focus on specific subsets
- **JSON Export**: Export your CSV data to JSON format for use in other applications
- **Data Correlation**: Calculate Pearson correlation coefficients between numeric columns
- **Anomaly Detection**: Find outliers in numeric data using Z-score method
- **Customizable Display**: Choose between light and dark mode for better visibility
