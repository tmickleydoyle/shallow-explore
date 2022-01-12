# shallow-explore
From the command line, quickly explore data from a CSV file.

`shallow-explore` is a [Golang](https://go.dev/) backed command-line tool for iterating over columns from a CSV file. This is a gut check tool to make sure the assumptions about the data are within the expected range of normal.

## How-To

After installation, run the following command to start analyzing data:

```bash
shallow-explore -path ~/Desktop/sample.csv
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
go get github.com/tmickleydoyle/shallow-explore
```

## Instructions for Installing Go

[Go docs](https://go.dev/)

### Installation with Homebrew

```bash
brew install go
```

## Why I Built This Tool

I find myself running and rerunning the same basic statistical analysis on data to get an understanding of how trends are moving. I figured why not make it easier and share it with everyone else! I hope this speeds up your decision making :heart:
