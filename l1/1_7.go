package main

// import (
// 	"fmt"
// 	"sync"
// )

// type SafeMap[K comparable, V interface{}] struct {
// 	mu        *sync.RWMutex
// 	container map[K]V
// }

// func NewSafeMap[K comparable, V interface{}]() *SafeMap[K, V] {
// 	mu := &sync.RWMutex{}
// 	container := make(map[K]V)
// 	return &SafeMap[K, V]{
// 		mu, container,
// 	}
// }

// func (sm *SafeMap[K, V]) Set(key K, val V) {
// 	sm.mu.Lock()
// 	defer sm.mu.Unlock()
// 	sm.container[key] = val
// }

// func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
// 	sm.mu.RLock()
// 	defer sm.mu.RUnlock()
// 	val, ok := sm.container[key]
// 	return val, ok
// }

// func main() {
// 	m := NewSafeMap[string, string]()
// 	var wg sync.WaitGroup

// 	for i := 0; i < 5; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 				key := fmt.Sprintf("key%d", i)
// 				val := fmt.Sprintf("value%d", i)
// 				m.Set(key, val)
// 				fmt.Printf("Set: %s, %s\n", key, val)
// 		}()
// 	}
	

// 	for i := 0; i < 10; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			key := fmt.Sprintf("key%d", i % 5)
// 			val, ok := m.Get(key)
// 			if ok {
// 				fmt.Printf("Get key: %s, val: %s\n", key, val)
// 			} else {
// 				fmt.Printf("No such key in map: %s\n", key) 
// 			}
// 		}()
// 	}

// 	wg.Wait()
// }
