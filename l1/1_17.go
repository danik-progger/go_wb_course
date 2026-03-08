package main

// import "fmt"

// func BinSearchInt(target int, arr []int) int {
// 	return BinSearch(target, arr, func(a, b int) bool { return a < b })
// }

// func BinSearchStr(target string, arr []string) int {
// 	return BinSearch(target, arr, func(a, b string) bool { return a < b })
// }

// func BinSearch[T comparable](target T, arr []T, compare func(l, r T) bool) int {
// 	l, r := 0, len(arr)-1

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
// 	res := BinSearchInt(4, slice)
// 	fmt.Println(res) // Should be 3
// 	res2 := BinSearchInt(10, slice)
// 	fmt.Println(res2) // Should be -1

// 	slice2 := []string{"a", "b", "c"}
// 	res3 := BinSearchStr("a", slice2)
// 	fmt.Println(res3) // Should be 0
// 	res4 := BinSearchStr("f", slice2)
// 	fmt.Println(res4) // Should be -1
// }
