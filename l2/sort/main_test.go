package main

import (
	"reflect"
	"testing"
)

func TestSorter(t *testing.T) {
	t.Run("BasicSort", func(t *testing.T) {
		s := &sorter{
			lines: []string{"c", "a", "b"},
		}
		s.sort()
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("ReverseSort", func(t *testing.T) {
		s := &sorter{
			lines:   []string{"c", "a", "b"},
			reverse: true,
		}
		s.sort()
		expected := []string{"c", "b", "a"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("NumericSort", func(t *testing.T) {
		s := &sorter{
			lines:   []string{"10", "2", "1"},
			numeric: true,
		}
		s.sort()
		expected := []string{"1", "2", "10"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("NumericReverseSort", func(t *testing.T) {
		s := &sorter{
			lines:   []string{"10", "2", "1"},
			numeric: true,
			reverse: true,
		}
		s.sort()
		expected := []string{"10", "2", "1"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("Unique", func(t *testing.T) {
		s := &sorter{
			lines: []string{"c", "a", "b", "a", "c"},
		}
		s.sort()
		s.lines = s.removeDuplicates()
		expected := []string{"a", "b", "c"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("ColumnSort", func(t *testing.T) {
		s := &sorter{
			lines:  []string{"a 2", "c 1", "b 3"},
			column: 2,
		}
		s.sort()
		expected := []string{"c 1", "a 2", "b 3"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("MonthSort", func(t *testing.T) {
		s := &sorter{
			lines: []string{"Mar", "Jan", "Feb"},
			month: true,
		}
		s.sort()
		expected := []string{"Jan", "Feb", "Mar"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("IgnoreBlanks", func(t *testing.T) {
		s := &sorter{
			lines:  []string{" c", " a", "b"},
			ignore: true,
		}
		s.sort()
		expected := []string{" a", "b", " c"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})

	t.Run("CheckSorted", func(t *testing.T) {
		s := &sorter{
			lines: []string{"a", "b", "c"},
		}
		if !s.isSorted() {
			t.Error("expected true, got false")
		}
	})

	t.Run("CheckNotSorted", func(t *testing.T) {
		s := &sorter{
			lines: []string{"c", "a", "b"},
		}
		if s.isSorted() {
			t.Error("expected false, got true")
		}
	})

	t.Run("HumanNumericSort", func(t *testing.T) {
		s := &sorter{
			lines: []string{"1G", "2K", "3M"},
			human: true,
		}
		s.sort()
		expected := []string{"2K", "3M", "1G"}
		if !reflect.DeepEqual(s.lines, expected) {
			t.Errorf("expected %v, got %v", expected, s.lines)
		}
	})
}
