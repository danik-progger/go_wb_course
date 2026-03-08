package main

// import (
// 	"fmt"
// 	"slices"
// )

// func intersect[T comparable](arr1 []T, arr2 []T) []T{
// 	res := []T{}

// 	for _, el := range arr1 {
// 		if slices.Contains(arr2, el) {
// 			res = append(res, el)
// 		}
// 	}
	
// 	return res
// }

// func main() {
// 	arr1 := []int{1, 2, 3, 4, 8, 11, 16}
// 	arr2 := []int{4, 3, 5, 6, 8, 16, -8}
// 	// Intersection is {3, 4, 8, 16}

// 	arr3 := []int8{1, 2, 3, 4, 8, 11, 16}
// 	arr4 := []int8{4, 3, 5, 6, 8, 16, -8}
// 	// Intersection is {3, 4, 8, 16}

// 	arr5 := []bool{true, false}
// 	arr6 := []bool{true, true, true}
// 	// Intersection is {true}
	
// 	arr7 := []string{"true", "abacaba", "false"}
// 	arr8 := []string{"true", "abacaba", "a"}
// 	// Intersection is {"true", "abacaba"}

// 	fmt.Println(intersect(arr1, arr2))
// 	fmt.Println(intersect(arr3, arr4))
// 	fmt.Println(intersect(arr5, arr6))
// 	fmt.Println(intersect(arr7, arr8))
// }
