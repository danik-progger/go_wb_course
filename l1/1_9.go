package main

// import "fmt"

// func writeToChannel(arr *[]int) <-chan int {
// 	out := make(chan int)
	
// 	go func() {
// 		for _, n := range *arr {
// 			out <- n
// 		}
// 		close(out)
// 	}()
	
// 	return out
// }

// func square(in <-chan int) <-chan int {
// 	out := make(chan int)
	
// 	go func() {
// 		for n := range in {
// 			out <- n * n
// 		}
// 		close(out)
// 	}()
	
// 	return out
// }

// func main() {
// 	arr := []int{1, 2, 3, 4, 5}
	
// 	// stage 1
// 	ch1 := writeToChannel(&arr)
	
// 	// stage 2
// 	ch2 := square(ch1)
	
// 	// stage 3
// 	for n := range ch2 {
// 		fmt.Printf("%d, ", n)
// 	}
// 	fmt.Println()
// }