package main

// import (
// 	"fmt"
// )

// func reverseString(s string) string {
// 	runes := []rune(s)

// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}

// 	return string(runes)
// }

// func main() {
// 	test := []string{
// 		"главрыба",
// 		"Hello, мир!",
// 		"🌟🚀💻",
// 		"Привет👋мир🌍",
// 		"тест",
// 	}

// 	for _, example := range test {
// 		fmt.Printf("%s → %s\n", example, reverseString(example))
// 	}
// }