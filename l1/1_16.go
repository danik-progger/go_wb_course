package main

// import "fmt"

// func QuickSort[T any](data []T, compare func(a, b T) bool) {
// 	if len(data) < 2 {
// 		return
// 	}
// 	l, r := 0, len(data)-1
// 	pivot := len(data) / 2
// 	data[pivot], data[r] = data[r], data[pivot]
// 	for i := range data {
// 		if compare(data[i], data[r]) {
// 			data[l], data[i] = data[i], data[l]
// 			l++
// 		}
// 	}
// 	data[l], data[r] = data[r], data[l]
// 	QuickSort(data[:l], compare)
// 	QuickSort(data[l+1:], compare)
// }

// type Cat struct {
// 	Name string
// 	Age  int
// }

// func main() {
// 	slice := []int{9, 4, 6, 2, 10, 3, 15, -5, 1}
// 	QuickSort(slice, func(a, b int) bool { return a < b })
// 	fmt.Println(slice)

// 	slice2 := []string{"f", "e", "d", "c", "b", "a"}
// 	QuickSort(slice2, func(a, b string) bool { return a < b })
// 	fmt.Println(slice2)

// 	slice3 := []Cat{{Name: "Harry", Age: 18}, {Name: "Tom", Age: 10}}
// 	QuickSort(slice3, func(a, b Cat) bool { return a.Age < b.Age })
// 	fmt.Println(slice3)
// }
