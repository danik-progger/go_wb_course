package main

// import (
// 	"fmt"
// 	"strings"
// )

// func reverseWords(s string, sep string) string {
// 	words := strings.Split(s, sep)
	
// 	for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
// 		words[i], words[j] = words[j], words[i]
// 	}
	
// 	res := strings.Join(words, sep)
	
// 	return res
// }

// func main() {
// 	test := []string{
// 		"Привет",
// 		"Тут уже много слов",
// 		"Я говорю как магистр йода",
// 	}

// 	for _, example := range test {
// 		fmt.Printf("%s → %s\n", example, reverseWords(example, " "))
// 	}
// }