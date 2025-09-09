package main

// import "fmt"

// func split(arr *[]float32) map[int][]float32 {
// 	res := make(map[int][]float32)
	
// 	for _, val := range *arr {
// 		key := int(val/10) * 10

// 		_, ok := res[key]
// 		if !ok {
// 			res[key] = make([]float32, 0)
// 		}
// 		res[key] = append(res[key], val)
// 	}

// 	return res
// }

// func main() {
// 	arr := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5}
// 	res := split(&arr)
	
// 	fmt.Println(res)
// }
