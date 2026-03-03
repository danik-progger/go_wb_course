package main

import (
	"fmt"
	"sort"
	"strings"
)

func FindAnograms(input []string) map[string][]string {
	res := make(map[string][]string)
	anograms := make(map[string][]string)
	for _, el := range input {
		lower := strings.ToLower(el)
		runes := []rune(lower)
		sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
		key := string(runes)
		anograms[key] = append(anograms[key], el)
	}

	for _, val := range anograms {
		if len(val) > 1 {
			sort.Strings(val)
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
