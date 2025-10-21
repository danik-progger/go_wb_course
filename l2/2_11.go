package main

import (
	"fmt"
	"sort"
)

func FindAnograms(input []string) map[string][]string {
	res := make(map[string][]string)
	anograms := make(map[string][]string)
	for _, el := range input {
		runes := []rune(el)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)
		if _, ok := anograms[key]; ok {
			anograms[key] = append(anograms[key], el)
		} else {
			anograms[key] = []string{el}
		}
	}

	for _, val := range anograms {
		if len(val) > 1 {
			res[val[0]] = val
		}
	}

	return res
}

func main() {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	result := FindAnograms(input)
	fmt.Println(result)
}
