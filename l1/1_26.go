package main

// import "fmt"

// func allLettersAppearOnce(s string) bool {
// 	set := make(map[rune]bool, 0) // Вобще хотелось бы set, но в go с ним тяжко
// 	for _, ch := range s {
// 		_, ok := set[ch]
// 		if ok {
// 			return false
// 		}
// 		set[ch] = true
// 	}

// 	return true
// }

// func main() {
// 	test1 := "abcdefg"
// 	test2 := "abacaba"

// 	fmt.Printf("%s has unique letters: %t\n", test1, allLettersAppearOnce(test1))
// 	fmt.Printf("%s has unique letters: %t\n", test2, allLettersAppearOnce(test2))
// }
