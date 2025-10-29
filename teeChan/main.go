package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func tee(ctx context.Context, in <-chan int, numChans int) []chan int {
	chans := make([]chan int, numChans)

	for i := range numChans {
		chans[i] = make(chan int)
	}

	go func() {
		for i := range numChans {
			defer close(chans[i])
		}

		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				wg := &sync.WaitGroup{}

				wg.Add(numChans)
				for i := range numChans {
					go func() {
						defer wg.Done()

						select {
						case <-ctx.Done():
							return
						case chans[i] <- val:
						}
					}()
				}

				wg.Wait()
			}
		}
	}()

	return chans
}

func generate() <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := range 10 {
			out <- i
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	numChans := 4
	chans := tee(ctx, generate(), numChans)

	wg := &sync.WaitGroup{}
	wg.Add(numChans)
	for i := range numChans {
		go func() {
			defer wg.Done()
			for val := range chans[i] {
				time.Sleep(time.Millisecond * 300)
				fmt.Println("Chan", i, val)
			}
		}()
	}

	wg.Wait()
}
