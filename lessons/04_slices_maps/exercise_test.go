package slicesmaps

import (
	"maps"
	"slices"
	"testing"
)

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3, 4}); got != 10 {
		t.Errorf("Sum([1 2 3 4]) = %d, want 10", got)
	}
	if got := Sum(nil); got != 0 {
		t.Errorf("Sum(nil) = %d, want 0", got)
	}
}

func TestMaxOf(t *testing.T) {
	if got, ok := MaxOf([]int{3, 9, 2}); got != 9 || !ok {
		t.Errorf("MaxOf([3 9 2]) = %d, %v, want 9, true", got, ok)
	}
	if got, ok := MaxOf([]int{-5, -2, -8}); got != -2 || !ok {
		t.Errorf("MaxOf([-5 -2 -8]) = %d, %v, want -2, true", got, ok)
	}
	if got, ok := MaxOf([]int{}); got != 0 || ok {
		t.Errorf("MaxOf([]) = %d, %v, want 0, false", got, ok)
	}
}

func TestEvens(t *testing.T) {
	got := Evens([]int{1, 2, 3, 4, 5, 6})
	if want := []int{2, 4, 6}; !slices.Equal(got, want) {
		t.Errorf("Evens([1..6]) = %v, want %v", got, want)
	}
	if got := Evens([]int{1, 3}); len(got) != 0 {
		t.Errorf("Evens([1 3]) = %v, want []", got)
	}
}

func TestReverse(t *testing.T) {
	in := []int{1, 2, 3}
	got := Reverse(in)
	if want := []int{3, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("Reverse([1 2 3]) = %v, want %v", got, want)
	}
	if want := []int{1, 2, 3}; !slices.Equal(in, want) {
		t.Errorf("Reverse modified its input: %v", in)
	}
}

func TestWordCount(t *testing.T) {
	got := WordCount("go is fun go")
	want := map[string]int{"go": 2, "is": 1, "fun": 1}
	if !maps.Equal(got, want) {
		t.Errorf("WordCount(%q) = %v, want %v", "go is fun go", got, want)
	}
}

func TestUnique(t *testing.T) {
	got := Unique([]string{"a", "b", "a", "c", "b"})
	if want := []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Errorf("Unique = %v, want %v", got, want)
	}
}
