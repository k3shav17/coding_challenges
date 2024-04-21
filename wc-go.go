package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

func main() {
	commandLineArgs := os.Args[1:]
	if commandLineArgs[0] == "-h" {
		Help()
		os.Exit(0)
	}

	noFlags(commandLineArgs)
	flag := commandLineArgs[0]
	filePath := commandLineArgs[1]
	if !fileExists(filePath) {
		fmt.Println("No such file or directory exists")
		os.Exit(-1)
	}
	fmt.Println(filePath)

	bytes, _ := readFile(filePath)
	Flags(flag, filePath, &bytes)
}

func Flags(flag, filePath string, content *[]byte) {
	switch flag {
	case "-c":
		NoOfBytes(content, filePath)
	case "-l":
		NoOfLines(content, filePath)
	case "-w":
		NoOfWords(content, filePath)
	case "-m":
		NoOfCharacters(content, filePath)
	default:
		Help()
	}
}

func Help() {
	fmt.Println("flags to be used with wc-go command")
	fmt.Println("\t-c \tThe number of bytes in each input file is written to the standard output")
	fmt.Println("\t-l \tThe number of lines in each input file is written to the standard output.")
	fmt.Println("\t-m \tThe number of characters in each input file is written to the standard output.",
		"\n\t\tIf the current locale does not support multibyte characters, this is equivalent to the -c option.",
		"\n\t\tThis will cancel out any prior usage of the -c option.")
	fmt.Println("\t-w\tThe number of words in each input file is written to the standard output.")
}

func NoOfBytes(content *[]byte, filePath string) {
	fmt.Print(len(string(*content)), "\t", filePath, "\n") // if it is -c
}

func NoOfLines(content *[]byte, filePath string) {
	fmt.Print(len(strings.Split(string(*content), "\n")), "\t", filePath, "\n") // if it is -l
}

func NoOfWords(content *[]byte, filePath string) {
	splitStr := strings.FieldsFunc(string(*content), func(c rune) bool {
		return unicode.IsSpace(c)
	})
	fmt.Print(len(splitStr), "\t", filePath, "\n") // if it is -w
}

func NoOfCharacters(content *[]byte, filePath string) {
	charCount := 0
	contentOfFile := string(*content)
	for _, r := range contentOfFile {
		if isChar(r) {
			charCount++
		}
	}
	fmt.Println(charCount, "\t", filePath, "\n") // if it is -m
}

func isChar(r rune) bool {
	return unicode.IsGraphic(r) || unicode.IsPrint(r) ||
		unicode.IsSpace(r) || unicode.IsSymbol(r) ||
		unicode.IsDigit(r) || unicode.IsLetter(r) ||
		unicode.IsMark(r) || unicode.IsNumber(r) ||
		unicode.IsPunct(r)
}

func readFile(filePath string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Unable to read the file %q", filePath)
	}
	return bytes, nil
}

func fileExists(filePath string) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		fmt.Printf("No such file or directory exists: %q", err)
		return false
	}
	return !fileInfo.IsDir()
}

func noFlags(commandLineArgs []string) {
	if len(commandLineArgs) == 1 && fileExists(commandLineArgs[0]) {
		var filePath string = commandLineArgs[0]
		content, _ := readFile(filePath)
		NoOfBytes(&content, filePath)
		NoOfLines(&content, filePath)
		NoOfWords(&content, filePath)
		NoOfCharacters(&content, filePath)
		os.Exit(0)
	}
}
