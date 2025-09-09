package main

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func worker(workerId int, ch <-chan int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for v := range ch {
// 		fmt.Printf("Worker %d: %d\n", workerId, v)
// 	}
// }

// func sender() (chan<- int, func()) {
// 	ch := make(chan int, 100)
// 	var wg sync.WaitGroup

// 	for i := 0; i < 5; i++ {
// 		wg.Add(1)
// 		go worker(i+1, ch, &wg)
// 	}

// 	cleanup := func() {
// 		close(ch)
// 		wg.Wait()
// 	}

// 	return ch, cleanup
// }

// func main() {
// 	fmt.Print("Enter number of seconds: ")
// 	var n int
// 	_, err := fmt.Scanf("%d", &n)
// 	if err != nil || n <= 0 {
// 		fmt.Println("🔴 Invalid number of workers, defaulting to 2")
// 		n = 2
// 	}

// 	ch, cleanup := sender()

// 	ctx, stop := context.WithTimeout(context.Background(), time.Second*time.Duration(n))
// 	defer stop()

// 	go func(n int) {
// 		i := 0
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			default:
// 				ch <- i
// 				i++
// 				time.Sleep(250 * time.Millisecond)
// 			}
// 		}
// 	}(n)

// 	<-ctx.Done()
// 	fmt.Println("\nshutting down ...")
// 	cleanup()
// 	fmt.Println("done")
// }
