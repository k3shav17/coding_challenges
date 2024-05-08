package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	Lines      string = "-l"
	Words      string = "-w"
	CharBytes  string = "-c"
	Characters string = "-m"
)

func readFromStdInput(option string) {
	reader := bufio.NewReader(os.Stdin)
	count := 0
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			count++
			break
		}
		if option == Lines {
			count++
		} else {
			count += Flags(option, &line)
		}
	}
	fmt.Println(count)
  os.Exit(0)
}

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

	bytes, _ := readFile(filePath)
	fmt.Println(Flags(flag, &bytes), filePath, "\t")
}

func Flags(flag string, content *[]byte) int {
	switch flag {
	case CharBytes:
		return NoOfBytes(content)
	case Lines:
		return NoOfLines(content)
	case Words:
		return NoOfWords(content)
	case Characters:
		return NoOfCharacters(content)
	default:
		Help()
	}
	return -1
}

func Help() {
	fmt.Println("flags to be used with wc-go command")
	fmt.Println("\t-c \tThe number of bytes in each input file is written to the standard output")
	fmt.Println("\t-l \tThe number of lines in each input file is written to the standard output.")
	fmt.Println("\t-m \tThe number of characters in each input file is written to the standard output.",
		"\n\t\tIf the current locale does not support multibyte characters, this is equivalent to the -c option.",
		"\n\t\tThis will cancel out any prior usage of the -c option.")
	fmt.Println("\t-w\tThe number of words in each input file is written to the standard output.")
  os.Exit(0)
}

func NoOfBytes(content *[]byte) int {
	return len(string(*content)) // if it is -c
}

func NoOfLines(content *[]byte) int {
	return len(strings.Split(string(*content), "\n")) // if it is -l
}

func NoOfWords(content *[]byte) int {
	splitStr := strings.FieldsFunc(string(*content), func(c rune) bool {
		return unicode.IsSpace(c)
	})
	return len(splitStr) // if it is -w
}

func NoOfCharacters(content *[]byte) int {
	charCount := 0
	contentOfFile := string(*content)
	for _, r := range contentOfFile {
		if isChar(r) {
			charCount++
		}
	}
	return charCount // if it is -m
}

func isChar(r rune) bool {
	return unicode.IsGraphic(r) || unicode.IsPrint(r) ||
		unicode.IsSpace(r) || unicode.IsSymbol(r) ||
		unicode.IsDigit(r) || unicode.IsLetter(r) ||
		unicode.IsMark(r) || unicode.IsNumber(r) ||
		unicode.IsPunct(r)
}

func readFile(filePath string) ([]byte, error) {
	bytes, err := os.ReadFile(filePath)
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
    fmt.Printf("%d %d %d %d %s\n", NoOfBytes(&content), NoOfLines(&content), NoOfWords(&content), NoOfCharacters(&content), filePath)
		os.Exit(0)
	}
}
