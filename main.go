package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/xuri/excelize/v2"
)

type ExitCode int

const (
	Matched = iota
	NotMatched
	ParameterError
)

var is_matched = false

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: xlgrep PATTERN FILENAME...")
		os.Exit(ParameterError)
	}

	reg, err := regexp.CompilePOSIX(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "PATTERN only support 'CompilePOSIX' syntax in go.")
		os.Exit(ParameterError)
	}

	files := os.Args[2:]

	for _, f := range files {
		grep(reg, f)
	}

	if is_matched {
		os.Exit(Matched)
	}
	os.Exit(NotMatched)
}

func grep(reg *regexp.Regexp, file string) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "'%s' open failed. xlgrep only support 'xlsx' format.\n", file)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "'%s' close failed.\n", file)
		}
	}()

	sheets := f.GetSheetList()
	for _, s := range sheets {
		rows, err := f.GetRows(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "'%s': row data read failed. skiped.\n", file)
			continue
		}

		for ridx, r := range rows {
			for cidx, c := range r {
				if reg.MatchString(c) {
					matched(file, s, ridx, cidx, c)
				}
			}
		}
	}
}

func matched(fileName string, sheetName string, rowIdx int, colIdx int, text string) {
	is_matched = true

	cname, err := excelize.CoordinatesToCellName(colIdx + 1, rowIdx + 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "'%s': cell coordinates translate failed.\n", fileName)
		return
	}
	fmt.Printf("%s:%s:%s:%s\n", fileName, sheetName, cname, text)
}
