package main

// import (
// 	"context"
// 	"fmt"
// 	"runtime"
// 	"sync"
// 	"time"
// )

// func exitWithContext(ctx context.Context) {
//  for {
//         select {
//         case <-ctx.Done():
//             return
//         default:
//             fmt.Println("Work...")
//             time.Sleep(1 * time.Second)
//         }
//     }
// }

// func exitsWithGoexit() {
	
// 	runtime.Goexit()
// }

	
// func exitWithChannel(done <-chan struct{}) {
//     for {
//         select {
//         case <-done:
//             return
//         default:
//             fmt.Println("Work...")
//             time.Sleep(1 * time.Second)
//         }
//     }
// }

// func waitTillCompletes(done chan<- bool) {
// 	fmt.Println("Worker: Starting...")
//     time.Sleep(2 * time.Second)
//     fmt.Println("Worker: Done.")
//     done <- true 
// }

// func exitWithPanic(wg *sync.WaitGroup) {
// 	defer wg.Done()
//     defer func() {
//         if r := recover(); r != nil {
//             fmt.Println("Finished in recovery after panic: \n", r)
//         }
//     }()

//     for {
//         fmt.Println("Work...")
//         time.Sleep(1 * time.Second)

//         if time.Now().Unix()%5 == 0 {
//             panic("AAAAA!")
//         }
//     }
// }


// func main() {
// 	// CONTEXT
// 	ctx, cancel := context.WithCancel(context.Background())
// 	go exitWithContext(ctx)
// 	time.Sleep(3 * time.Second)
// 	cancel()
// 	fmt.Println("Stopped with context\n")

// 	// NOTIFICATION CHANNEL
//     done := make(chan struct{})
// 	go exitWithChannel(done)
// 	time.Sleep(3 * time.Second)
// 	close(done)
// 	fmt.Println("Stopped with notification channel\n")
	
	
// 	// WAIT TILL COMPLETES
// 	ch := make(chan bool)
// 	go waitTillCompletes(ch)
// 	<- ch
// 	fmt.Println("Waited with channel till goroutine finished\n")
	
// 	// BY CONDITION
// 	running := true
// 	var wg sync.WaitGroup
// 	wg.Add(1)
	
// 	go func() {	
// 		defer wg.Done()
//   		for running {
//             fmt.Println("Work...")
//             time.Sleep(500 * time.Millisecond)
//         }
// 	}()
// 	time.Sleep(3 * time.Second)
// 	running = false
// 	wg.Wait()
// 	fmt.Println("Stopped by condition\n")
	
// 	running = true
// 	go func() {
// 		for {
// 			if !running {
// 				runtime.Goexit()
// 			}
			
// 			fmt.Println("Work...")
//            	time.Sleep(500 * time.Millisecond)
// 		}
// 	}()
// 	time.Sleep(3 * time.Second)
// 	running = false
// 	fmt.Println("Stopped with Goexit\n")
	
// 	// PANIC AND RECOVER
// 	wg.Add(1)
// 	go exitWithPanic(&wg)
// 	wg.Wait()
	
// 	// STOPED BY FINISHING MAIN GOROUTINE
// 	go func() {
// 		for {
// 			fmt.Println("Work...")
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}()
// 	time.Sleep(3 * time.Second)
// 	fmt.Println("Main goroutine finishes so this goroutine will be cut off\n")
// }
