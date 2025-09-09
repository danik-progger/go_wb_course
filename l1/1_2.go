package main

// import (
// 	"fmt"
// 	"sync"
// )

// func main() {
// 	arr := [5]int{2, 4, 6, 8, 10}
// 	var wg sync.WaitGroup
// 	wg.Add(len(arr))

// 	for i := range arr {
// 		go func(idx int) {
// 			defer wg.Done()
// 			arr[idx] = arr[idx] * arr[idx]
// 		}(i)
// 	}

// 	wg.Wait()

// 	for _, n := range arr {
// 		fmt.Println(n)
// 	}
// }
