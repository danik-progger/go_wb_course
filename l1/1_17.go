package main

// import "fmt"


// func BinSearch[T comparable](target T, arr []T, compare func (l, r T) bool) int {
// 	l, r := 0, len(arr) - 1
	
// 	for l <= r {
// 		m := (l + r) / 2
// 		if arr[m] == target {
// 			return m
// 		}
// 		if compare(arr[m], target) {
// 			l = m + 1
// 		} else {
// 			r = m - 1
// 		}
// 	} 

// 	return -1
// }

// func main() {
// 	slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
//     res := BinSearch(4, slice, func(a, b int) bool { return a < b })
//     fmt.Println(res) // Should be 3
//     res2 := BinSearch(10, slice, func(a, b int) bool { return a < b })
//     fmt.Println(res2) // Should be -1
// }
