package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func fanIn(ctx context.Context, chans []chan int) chan int {
	out := make(chan int)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(out)
		for _, ch := range chans {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case v, ok := <-ch:
						if !ok {
							return
						}

						select {
						case out <- v:
						case <-ctx.Done():
							return
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		wg.Wait()
	}()

	return out
}

func fanOut(in chan int, numChans int, f func(int) int) []chan int {
	chans := make([]chan int, numChans)

	for i := range numChans {
		chans[i] = pipeline(in, f)
	}

	return chans
}

func pipeline(in chan int, f func(int) int) chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for v := range in {
			out <- f(v)
		}
	}()

	return out
}

func generate() chan int {
	in := make(chan int)

	go func() {
		defer close(in)
		for i := range 100000 {
			in <- i
		}
	}()

	return in
}

func pool(in chan int, numWorkers int, f func(int) int) chan int {
	out := make(chan int)
	wg := &sync.WaitGroup{}

	go func() {
		defer close(out)
		for range numWorkers {
			wg.Add(1)
			go worker(in, out, f, wg)
		}
		wg.Wait()
	}()

	return out
}

func worker(in chan int, out chan int, f func(int) int, wg *sync.WaitGroup) {
	defer wg.Done()

	for v := range in {
		out <- f(v)
	}
}

func cpuIntensive() {
	c := 0

	for range 10000 {
		c++
	}
}

func ioIntensive() {
	time.Sleep(100 * time.Microsecond)
}

func compute(a int) int {
	cpuIntensive()
	ioIntensive()
	return a * a
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numWorkeers := 500

	fmt.Println("Fan")
	start := time.Now()
	for v := range fanIn(ctx, fanOut(generate(), numWorkeers, compute)) {
		fmt.Println(v)
	}

	timeFanIn := time.Since(start)

	fmt.Println("Pool")
	start2 := time.Now()
	for v := range pool(generate(), numWorkeers, compute) {
		fmt.Println(v)
	}

	timePool := time.Since(start2)

	fmt.Println("time fanin:", timeFanIn)
	fmt.Println("time pool:", timePool)
}
