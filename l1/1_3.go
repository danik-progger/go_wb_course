package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"time"
// )

// func worker(workerId int, ch <-chan int, wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	for v := range ch {
// 		fmt.Printf("Worker %d: %d\n", workerId, v)
// 	}
// }

// func sender(numWorkers int) (chan<- int, func()) {
// 	ch := make(chan int, 100)
// 	var wg sync.WaitGroup

// 	for i := 0; i < numWorkers; i++ {
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
// 	fmt.Print("Enter number of workers: ")
// 	var n int
// 	_, err := fmt.Scanf("%d", &n)
// 	if err != nil || n <= 0 {
// 		fmt.Println("🔴 Invalid number of workers, defaulting to 2")
// 		n = 2
// 	}

// 	ch, cleanup := sender(n)

// 	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
// 	defer stop()

// 	go func() {
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
// 	}()

// 	<-ctx.Done()
// 	fmt.Println("\nshutting down ...")
// 	cleanup()
// 	fmt.Println("done")
// }
