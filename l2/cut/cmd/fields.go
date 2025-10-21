package cut

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Fields struct {
	indices []int
}

func NewFields(fieldsStr string) (*Fields, error) {
	fieldSet := make(map[int]struct{})
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}
			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", rangeParts[0])
			}
			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", rangeParts[1])
			}
			if start <= 0 || end <= 0 || start > end {
				return nil, fmt.Errorf("invalid range: %d-%d", start, end)
			}
			for i := start; i <= end; i++ {
				fieldSet[i-1] = struct{}{} // 1-based to 0-based
			}
		} else {
			field, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			if field <= 0 {
				return nil, fmt.Errorf("field numbers must be positive: %d", field)
			}
			fieldSet[field-1] = struct{}{} // 1-based to 0-based
		}
	}

	var indices []int
	for field := range fieldSet {
		indices = append(indices, field)
	}
	sort.Ints(indices)
	return &Fields{indices: indices}, nil
}

func (f *Fields) Extract(parts []string) []string {
	var resultParts []string
	for _, fieldIndex := range f.indices {
		if fieldIndex < len(parts) {
			resultParts = append(resultParts, parts[fieldIndex])
		}
	}
	return resultParts
}
