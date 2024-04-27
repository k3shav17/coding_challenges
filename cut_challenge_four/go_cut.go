package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var delimiter = "\t"

func main() {

	var dataForLines []string
	commandLineArgs := os.Args[1:]
	flag := commandLineArgs[0]

	splitFunc := func(c rune) bool {
		return c == ' ' || c == ','
	}
	dataForLines = strings.FieldsFunc(flag[2:], splitFunc)
	var cols []int
	for _, val := range dataForLines {
		col, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Please enter a column or line to cut")
			os.Exit(0)
		}
		cols = append(cols, col)
	}

	filePath := commandLineArgs[1] // if there is no delimiter

	if len(commandLineArgs) >= 3 { // if there is a delimiter
		filePath = commandLineArgs[2]
		delimiter = commandLineArgs[1][len(commandLineArgs[1])-1:]
	}

	content, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	processFile(content, cols)
	defer content.Close()
}

func processFile(fileContent *os.File, cols []int) {
	scanner := bufio.NewScanner(fileContent)
	for scanner.Scan() {
		line := scanner.Text()
		doCutFields(line, delimiter, cols)
	}
}

func doCutFields(line, delimiter string, fields []int) {
	tokens := strings.Split(line, delimiter)
	for _, col := range fields {
		if col <= len(tokens) {
			fmt.Print(tokens[col-1], " ")
		}
	}
	fmt.Println()
}
