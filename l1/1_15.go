package main

// import (
// 	"fmt"
// 	"strings"
// )

// func createHugeString(size int) string {
// 	return strings.Repeat("A", size)
// }

// func someFunc() string {
// 	v := createHugeString(10_000)
	
// 	if len(v) < 100 {
// 		return v
// 	}
// 	justString := string([]byte(v[:100]))
// 	return justString
// }

// func main() {
// 	res := someFunc()
// 	fmt.Printf("Result: %s\n", res)
// 	fmt.Printf("Length: %d\n", len(res))
// }