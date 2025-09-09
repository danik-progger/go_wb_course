package main

// import (
// 	"fmt"
// 	"math"
// 	"math/big"
// )

// func main() {
// 	// Примеры чисел больше 2^20 (больше 1 миллиона)
// 	a := int64(2000000) // 2 * 10^6
// 	b := int64(3000000) // 3 * 10^6

// 	fmt.Println("Программа для арифметических операций с большими числами")
// 	fmt.Printf("Числа для операций: a = %d, b = %d\n", a, b)
// 	fmt.Println("\n1. Использование встроенного типа int64:")

// 	// Арифметические операции с использованием int64
// 	fmt.Printf("a + b = %d\n", a+b)
// 	fmt.Printf("a - b = %d\n", a-b)
// 	fmt.Printf("a * b = %d\n", a*b)
// 	fmt.Printf("a / b = %.2f\n", float64(a)/float64(b))

// 	// Демонстрация возможного переполнения при использовании int64
// 	maxInt64 := int64(math.MaxInt64)
// 	fmt.Printf("\nМаксимальное значение int64: %d\n", maxInt64)

// 	bigA := int64(9000000000000000000) // 9 * 10^18, приближается к максимуму int64
// 	bigB := int64(2)

// 	fmt.Printf("bigA = %d, bigB = %d\n", bigA, bigB)
// 	fmt.Printf("bigA * bigB = ")
// 	if bigA > maxInt64/bigB {
// 		fmt.Println("Ошибка: переполнение int64")
// 	} else {
// 		fmt.Printf("%d\n", bigA*bigB)
// 	}

// 	// Использование math/big для работы с произвольно большими числами
// 	fmt.Println("\n2. Использование math/big для произвольно больших чисел:")

// 	// Создание big.Int из строк для предотвращения переполнения
// 	bigIntA := new(big.Int)
// 	bigIntA.SetString("9000000000000000000000000000000", 10) // 9 * 10^30

// 	bigIntB := new(big.Int)
// 	bigIntB.SetString("5000000000000000000000000000000", 10) // 5 * 10^30

// 	fmt.Printf("bigIntA = %s\n", bigIntA.String())
// 	fmt.Printf("bigIntB = %s\n", bigIntB.String())

// 	// Сложение
// 	sum := new(big.Int).Add(bigIntA, bigIntB)
// 	fmt.Printf("bigIntA + bigIntB = %s\n", sum.String())

// 	// Вычитание
// 	diff := new(big.Int).Sub(bigIntA, bigIntB)
// 	fmt.Printf("bigIntA - bigIntB = %s\n", diff.String())

// 	// Умножение
// 	product := new(big.Int).Mul(bigIntA, bigIntB)
// 	fmt.Printf("bigIntA * bigIntB = %s\n", product.String())

// 	// Деление
// 	quotient := new(big.Rat).SetFrac(bigIntA, bigIntB)
// 	fmt.Printf("bigIntA / bigIntB = %s\n", quotient.FloatString(10))
// }