package main

import (
	"reflect"
	"testing"
)

func TestFindAnograms(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name:  "basic anagrams",
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{
			name:  "case insensitive",
			input: []string{"Пятак", "ПЯТКА", "тяпка"},
			expected: map[string][]string{
				"ПЯТКА": {"ПЯТКА", "Пятак", "тяпка"},
			},
		},
		{
			name:     "no anagrams",
			input:    []string{"стол", "стул", "шкаф"},
			expected: map[string][]string{},
		},
		{
			name:  "sorted output",
			input: []string{"пятка", "пятак", "тяпка"},
			expected: map[string][]string{
				"пятак": {"пятак", "пятка", "тяпка"},
			},
		},
		{
			name:     "empty input",
			input:    []string{},
			expected: map[string][]string{},
		},
		{
			name:  "single word",
			input: []string{"слово"},
			expected: map[string][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindAnograms(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FindAnograms(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
