package main

// import (
// 	"fmt"
// 	"sync/atomic"
// )

// type Counter struct {
// 	val int32
// }

// func (c *Counter) Increment() {
// 	atomic.AddInt32(&c.val, 1)
// }

// func (c *Counter) Get() int32 {
// 	return c.val
// }

// func main() {
// 	c := Counter{0}

// 	for _ = range 1000 {
// 		c.Increment()
// 	}
	
// 	fmt.Println(c.Get())
// }
