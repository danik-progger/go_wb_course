package main

// import (
// 	"fmt"
// 	"slices"
// )

// type MySet[T comparable] struct {
// 	container []T
// }

// func (s *MySet[T]) Has(el T) bool {
// 	return slices.Contains(s.container, el)
// }

// func (s *MySet[T]) Add(el T) {
// 	if !s.Has(el) {
// 		s.container = append(s.container, el)
// 	}
// }

// func (s *MySet[T]) Print() {
// 	fmt.Println(s.container)
// }

// func SetFromSlice[T comparable](s *[]T) *MySet[T] {
// 	var	set MySet[T]
// 	for _, el := range *s {
// 		set.Add(el)
// 	}
	
// 	return &set
// }


// func main () {
// 	arr := []string{"cat", "cat", "dog", "cat", "tree"}
// 	set := SetFromSlice(&arr)
// 	set.Print()
// }