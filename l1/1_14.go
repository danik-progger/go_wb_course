package main

// import "fmt"

// func PrintType[T interface{}](v T) {
// 	switch any(v).(type) {
// 	case int:
// 		fmt.Println("int")
// 	case string:
// 		fmt.Println("string")
// 	case bool:
// 		fmt.Println("bool")
// 	default:
//         fmt.Printf("%T\n", v)
// 	}
// }

// func main() {
// 	var ch chan int
// 	test := []interface{}{"Code", true, 5, ch}
// 	for _, el := range test {
// 		PrintType(el)
// 	}
// }