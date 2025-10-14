package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func StringUnpack(s string) (string, error) {
	var result strings.Builder
	input := []rune(s)
	i := 0

	for i < len(input) {
		currChar := input[i]

		if currChar == '\\' {
			i++
			if i >= len(input) {
				return "", errors.New("invalid string: ends with an escape character")
			}
			currChar = input[i]
		} else if unicode.IsDigit(currChar) {
			return "", errors.New("invalid string: character to be repeated cannot be a digit")
		}

		i++

		startOfNum := i
		for i < len(input) && unicode.IsDigit(input[i]) {
			i++
		}
		repeatCount := 1
		if startOfNum < i {
			numStr := string(input[startOfNum:i])
			repeatCount, _ = strconv.Atoi(numStr)
		}

		for range repeatCount {
			result.WriteRune(currChar)
		}
	}

	return result.String(), nil
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter string: ")
	text, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("🔴 Failed to read from console", err)
		return
	}
	text = strings.TrimSpace(text)

	unpacked, err := StringUnpack(text)
	if err != nil {
		fmt.Printf("🔴  Error: %v\n", err)
	} else {
		fmt.Printf("🟢 Unpacked string: %s\n", unpacked)
	}
}
