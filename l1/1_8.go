package main

// import (
// 	"fmt"
// )

// func setBit(n int64, pos int64, value bool) int64 {
// 	mask := int64(1) << (pos - 1)
// 	if value {
// 		return n | mask
// 	}

// 	return n &^ mask
// }

// func main() {
// 	var n, pos, val int64
// 	fmt.Print("Enter number: ")
// 	_, err := fmt.Scanf("%d", &n)
// 	if err != nil {
// 		fmt.Println("🔴 Error reading number")
// 	}
// 	fmt.Print("Enter bit position: ")
// 	_, err = fmt.Scanf("%d", &pos)
// 	if err != nil {
// 		fmt.Println("🔴 Error reading bit position")
// 	}
// 	fmt.Print("Enter bit value: ")
// 	_, err = fmt.Scanf("%d", &val)
// 	if err != nil {
// 		fmt.Println("🔴 Error reading bit value")
// 	}

// 	res := setBit(n, pos, val == 1)
// 	fmt.Println(res)
// }
